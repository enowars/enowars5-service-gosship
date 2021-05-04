package chat

import (
	"fmt"
	"gosship/pkg/database"
	"io"
	"sync"

	"github.com/gliderlabs/ssh"
	"github.com/kyokomi/emoji/v2"
	"github.com/logrusorgru/aurora/v3"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const title = "              _____ _____ _    _ _       \n             / ____/ ____| |  | (_)      \n   __ _  ___| (___| (___ | |__| |_ _ __  \n  / _` |/ _ \\\\___ \\\\___ \\|  __  | | '_ \\ \n | (_| | (_) |___) |___) | |  | | | |_) |\n  \\__, |\\___/_____/_____/|_|  |_|_| .__/ \n   __/ |                          | |    \n  |___/                           |_|    "

type Host struct {
	Log      *logrus.Logger
	mu       sync.RWMutex
	users    map[string]*User
	msgChan  chan Message
	Database *database.Database
}

func (h *Host) HandleNewSession(session ssh.Session) {
	code := 0
	err := h.handleNewSessionWithError(session)
	if err != nil {
		h.Log.Error(err)
		code = -1
	}
	if err := session.Exit(code); err != nil && err != io.EOF {
		h.Log.Error(err)
	}
}

func (h *Host) AddUser(u *User) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.users[u.Name]; ok {
		return false
	}
	h.users[u.Name] = u
	h.Log.Printf("[%s] added\n", u.Name)
	return true
}

func (h *Host) RemoveUser(name string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.users, name)
	h.Log.Printf("[%s] removed\n", name)
}

func (h *Host) writeLineLogError(w io.Writer, s string) {
	_, err := io.WriteString(w, s+"\n")
	if err != nil {
		h.Log.Error(err)
	}
}

func (h *Host) handleNewSessionWithError(session ssh.Session) error {
	_, _, isPty := session.Pty()
	if !isPty {
		h.writeLineLogError(session, "No PTY requested!")
		return fmt.Errorf("no PTY requested")
	}

	u, err := NewUser(h.Database, session)
	if err != nil {
		if err == ErrFingerprintDoesNotMatch {
			h.writeLineLogError(session, "The provided public key does not match with the username!")
		}
		if err == ErrFingerprintAlreadyRegistered {
			h.writeLineLogError(session, "The provided public key is already linked to a username!")
		}
		return err
	}

	h.Log.Printf("[%s] new session: fingerprint=(%s)\n", u.Name, u.Fingerprint)
	if !h.AddUser(u) {
		err := u.WriteLine(aurora.Sprintf("%s is already logged in!", aurora.Red(u.Name)))
		if err != nil {
			return err
		}
		return fmt.Errorf("[%s] already logged in", u.Name)
	}
	defer h.RemoveUser(u.Name)
	err = u.WriteLine(emoji.Sprintf("%s\n\n:unicorn:Welcome %s! You are now in room %s.\n", aurora.Green(title), aurora.Magenta(u.Name), aurora.Yellow(u.CurrentRoom)))
	if err != nil {
		return err
	}
	oldMessages, err := h.Database.GetRecentMessagesForUserAndRoom(u.Id, u.CurrentRoom)
	if err != nil {
		return err
	}
	for _, oldMsg := range oldMessages {
		// skip dm history
		if oldMsg.Type == database.MessageType_DIRECT {
			continue
		}
		conMsg, err := h.ConvertMessageEntryToMessage(oldMsg)
		if err != nil {
			h.Log.Error(err)
			continue
		}
		err = u.WriteMessage(conMsg)
		if err != nil {
			h.Log.Error(err)
		}
	}
	h.RoomAnnouncement(u.CurrentRoom, aurora.Sprintf("%s joined the room.", u.RenderName()))
	for {
		line, err := u.Term.ReadLine()
		if err != nil {
			if err != io.EOF {
				h.Log.Error(err)
			}
			break
		}
		if line == "" {
			_, _ = u.Term.Write([]byte{})
			continue
		}
		parsedMessage, err := ParseMessage(line, u)
		if err != nil {
			_ = u.WriteLine(aurora.Sprintf(aurora.Red("error: %s"), err.Error()))
			continue
		}
		h.RouteMessage(parsedMessage)
	}
	h.RoomAnnouncement(u.CurrentRoom, aurora.Sprintf("%s left the room.", u.RenderName()))
	return nil
}

func (h *Host) HandlePublicKey(ctx ssh.Context, key ssh.PublicKey) bool {
	h.Log.Printf("new connection (%s) with key type: %s\n", ctx.RemoteAddr(), key.Type())
	return true
}

func (h *Host) Serve() {
	for msg := range h.msgChan {
		var msgEntry database.MessageEntry
		skipSave := false
		h.Log.Println(msg.String())
		switch v := msg.(type) {
		case *PublicMessage:
			msgEntry.Type = database.MessageType_PUBLIC
			msgEntry.Body = v.Body
			msgEntry.Timestamp = timestamppb.New(v.Timestamp)
			msgEntry.From = v.From.Id
			msgEntry.Room = v.Room
			h.sendMessageToAllUsersInRoom(v)
		case *RoomAnnouncementMessage:
			msgEntry.Type = database.MessageType_ROOM_ANNOUNCEMENT
			msgEntry.Body = v.Body
			msgEntry.Timestamp = timestamppb.New(v.Timestamp)
			msgEntry.Room = v.Room
			h.sendMessageToAllUsersInRoom(v)
		case *AnnouncementMessage:
			msgEntry.Type = database.MessageType_ANNOUNCEMENT
			msgEntry.Body = v.Body
			msgEntry.Timestamp = timestamppb.New(v.Timestamp)
			h.sendMessageToAllUsers(v)
		case *DirectMessage:
			toId := h.resolveUserNameToID(v.To)
			msgEntry.Type = database.MessageType_DIRECT
			msgEntry.Body = v.Body
			msgEntry.Timestamp = timestamppb.New(v.Timestamp)
			msgEntry.From = v.From.Id
			msgEntry.To = toId
			skipSave = toId == ""
			h.sendMessageToUser(v)
		case *CommandMessage:
			skipSave = true
			h.handleUserCommand(v)
		default:
			skipSave = true
			h.Log.Error("unknown message type")
		}
		if skipSave {
			continue
		}
		if err := h.Database.AddMessageEntry(&msgEntry); err != nil {
			h.Log.Error(err)
		}
	}
}

func (h *Host) sendMessageToAllUsersInRoom(msg Message) {
	var room string
	if v, ok := msg.(*PublicMessage); ok {
		room = v.Room
	}
	if v, ok := msg.(*RoomAnnouncementMessage); ok {
		room = v.Room
	}
	if room == "" {
		h.Log.Error("room not found in message")
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, u := range h.users {
		if u.CurrentRoom != room {
			continue
		}
		err := u.WriteMessage(msg)
		if err != nil {
			h.Log.Error(err)
		}
	}
}

func (h *Host) sendMessageToAllUsers(msg *AnnouncementMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, u := range h.users {
		err := u.WriteMessage(msg)
		if err != nil {
			h.Log.Error(err)
		}
	}
}

func (h *Host) resolveUserNameToID(name string) string {
	//TODO: use database to resolve username
	h.mu.RLock()
	defer h.mu.RUnlock()
	if u, ok := h.users[name]; ok {
		return u.Id
	}
	return ""
}

func (h *Host) ConvertMessageEntryToMessage(me *database.MessageEntry) (Message, error) {
	rm := &rawMessage{
		Timestamp: me.Timestamp.AsTime().Local(),
		Body:      me.Body,
	}
	var from *User
	if me.From != "" {
		fromEntry, err := h.Database.GetUserById(me.From)
		if err != nil {
			return nil, err
		}
		from = &User{
			Id:    me.From,
			Name:  fromEntry.Name,
			Dummy: true,
		}
	}
	var to *User
	if me.To != "" {
		toEntry, err := h.Database.GetUserById(me.To)
		if err != nil {
			return nil, err
		}
		to = &User{
			Id:    me.To,
			Name:  toEntry.Name,
			Dummy: true,
		}
	}

	switch me.Type {
	case database.MessageType_PUBLIC:
		return &PublicMessage{
			rawMessage: rm,
			From:       from,
			Room:       me.Room,
		}, nil
	case database.MessageType_DIRECT:
		if to == nil {
			return nil, fmt.Errorf("recipient not found")
		}
		return &DirectMessage{
			rawMessage: rm,
			From:       from,
			To:         to.Name,
			ToResolved: to,
		}, nil
	case database.MessageType_ROOM_ANNOUNCEMENT:
		return &RoomAnnouncementMessage{
			rawMessage: rm,
			Room:       me.Room,
		}, nil
	case database.MessageType_ANNOUNCEMENT:
		return &AnnouncementMessage{
			rawMessage: rm,
		}, nil
	}
	return nil, fmt.Errorf("invalid message type")
}

func (h *Host) sendMessageToUser(msg *DirectMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	to, ok := h.users[msg.To]
	if !ok {
		err := msg.From.WriteLine(aurora.Sprintf(aurora.Yellow("user %s not found on the server."), aurora.Red(msg.To)))
		if err != nil {
			h.Log.Error(err)
		}
		return
	}
	for _, u := range []*User{msg.From, to} {
		err := u.WriteMessage(msg)
		if err != nil {
			h.Log.Error(err)
		}
	}
	to.LastDmRecipient = msg.From.Name
}

func (h *Host) handleUserCommand(msg *CommandMessage) {
	cmd := FindCommand(msg.Cmd)
	if cmd == nil {
		err := msg.From.WriteLine(aurora.Sprintf(aurora.Yellow("command %s not found. use /help to list all available commands."), aurora.Red(msg.Cmd)))
		if err != nil {
			h.Log.Error(err)
		}
		return
	}
	err := cmd.Handler(h, msg)
	if err != nil {
		_ = msg.From.WriteLine(aurora.Sprintf(aurora.Red("command error: %s"), err.Error()))
	}
}

func (h *Host) Announcement(msg string) {
	h.msgChan <- NewAnnouncementMessage(msg)

}

func (h *Host) RoomAnnouncement(room, msg string) {
	h.msgChan <- NewRoomAnnouncementMessage(room, msg)
}

func (h *Host) RouteMessage(msg Message) {
	h.msgChan <- msg
}

func NewHost(log *logrus.Logger, db *database.Database) *Host {
	return &Host{
		Log:      log,
		users:    make(map[string]*User),
		msgChan:  make(chan Message, 10),
		Database: db,
	}
}
