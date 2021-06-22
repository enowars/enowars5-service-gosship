package chat

import (
	"fmt"
	"gosship/pkg/database"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/kyokomi/emoji/v2"
	"github.com/logrusorgru/aurora/v3"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const title = "              _____ _____ _    _ _       \n             / ____/ ____| |  | (_)      \n   __ _  ___| (___| (___ | |__| |_ _ __  \n  / _` |/ _ \\\\___ \\\\___ \\|  __  | | '_ \\ \n | (_| | (_) |___) |___) | |  | | | |_) |\n  \\__, |\\___/_____/_____/|_|  |_|_| .__/ \n   __/ |                          | |    \n  |___/                           |_|    "

type Host struct {
	Log               *logrus.Logger
	usersMu           sync.RWMutex
	users             map[string]*User
	msgChan           chan Message
	Database          *database.Database
	roomsMu           sync.RWMutex
	Rooms             map[string]*database.RoomEntry
	cleanupTickerStop chan struct{}

	// hide user joined/left messages
	DisableUserAnnouncements bool

	DisableRoomAnnouncements   bool
	DisableServerAnnouncements bool
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
	h.usersMu.Lock()
	defer h.usersMu.Unlock()
	if _, ok := h.users[u.Name]; ok {
		return false
	}
	h.users[u.Name] = u
	h.Log.Printf("[%s] added", u.Name)
	return true
}

func (h *Host) RemoveUser(name string) {
	h.usersMu.Lock()
	defer h.usersMu.Unlock()
	delete(h.users, name)
	h.Log.Printf("[%s] removed", name)
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

	h.Log.Printf("[%s] new session: fingerprint=%s", u.Name, u.Fingerprint)
	if !h.AddUser(u) {
		err := u.WriteLine(aurora.Sprintf("%s is already logged in!", aurora.Red(u.Name)))
		if err != nil {
			return err
		}
		return fmt.Errorf("[%s] already logged in", u.Name)
	}
	defer h.RemoveUser(u.Name)

	if !h.HasRoom(u.CurrentRoom) {
		if err := u.UpdateCurrentRoom("default"); err != nil {
			return err
		}
	}

	var welcomeMessage strings.Builder
	welcomeMessage.WriteString(aurora.Sprintf("%s\n\n", aurora.Green(title)))
	welcomeMessage.WriteString(aurora.Sprintf(":unicorn:welcome %s!\n", u.RenderName(true)))
	welcomeMessage.WriteString(aurora.Sprintf("you are now in room %s and %s\n", aurora.Blue(u.CurrentRoom), h.ServerInfo()))
	welcomeMessage.WriteString(aurora.Sprintf("use %s to list all available commands.\n", aurora.Green("/help")))
	err = u.WriteLine(emoji.Sprint(welcomeMessage.String()))
	if err != nil {
		return err
	}

	err = h.ShowRecentMessages(u, false)
	if err != nil {
		return err
	}

	h.JoinRoomAnnouncement(u)
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
	h.LeftRoomAnnouncement(u)
	return nil
}

func (h *Host) HandlePublicKey(ctx ssh.Context, key ssh.PublicKey) bool {
	h.Log.Printf("new connection (%s) with key type: %s", ctx.RemoteAddr(), key.Type())
	return true
}

func (h *Host) Serve() {
	for msg := range h.msgChan {
		var msgEntry database.MessageEntry
		saveToDatabase := false
		h.Log.Debug(msg.String())
		switch v := msg.(type) {
		case *PublicMessage:
			saveToDatabase = true
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
			saveToDatabase = toId != ""
			h.sendMessageToUser(v, toId)
		case *CommandMessage:
			go h.handleUserCommand(v)
		default:
			h.Log.Error("unknown message type")
		}
		if !saveToDatabase {
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
	h.usersMu.RLock()
	defer h.usersMu.RUnlock()
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
	h.usersMu.RLock()
	defer h.usersMu.RUnlock()
	for _, u := range h.users {
		err := u.WriteMessage(msg)
		if err != nil {
			h.Log.Error(err)
		}
	}
}

func (h *Host) resolveUserNameToID(name string) string {
	h.usersMu.RLock()
	defer h.usersMu.RUnlock()
	if u, ok := h.users[name]; ok {
		return u.Id
	}
	id, _, err := h.Database.FindUserByPredicate(func(entry *database.UserEntry) bool {
		return entry.Name == name
	})
	if err != nil {
		return ""
	}
	return id
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

func (h *Host) sendMessageToUser(msg *DirectMessage, toId string) {
	h.usersMu.RLock()
	defer h.usersMu.RUnlock()
	to, ok := h.users[msg.To]
	if !ok && toId == "" {
		err := msg.From.WriteLine(aurora.Sprintf(aurora.Yellow("user %s not found on the server."), aurora.Red(msg.To)))
		if err != nil {
			h.Log.Error(err)
		}
		return
	}
	recipients := []*User{msg.From}
	if ok {
		recipients = append(recipients, to)
		to.LastDmRecipient = msg.From.Name
	} else if toId != "" {
		err := msg.From.WriteLine(aurora.Sprintf(aurora.Yellow("user %s is currently not online. the message still was sent and can be retrieved with the /history command."), aurora.Red(msg.To)))
		if err != nil {
			h.Log.Error(err)
		}
	}

	for _, u := range recipients {
		err := u.WriteMessage(msg)
		if err != nil {
			h.Log.Error(err)
		}
	}
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
	if h.DisableServerAnnouncements {
		return
	}
	h.RouteMessage(NewAnnouncementMessage(msg))
}

func (h *Host) RoomAnnouncement(room, msg string) {
	if h.DisableRoomAnnouncements {
		return
	}
	h.RouteMessage(NewRoomAnnouncementMessage(room, msg))
}

func (h *Host) RouteMessage(msg Message) {
	select {
	case h.msgChan <- msg:
	case <-time.After(time.Second * 5):
		h.Log.Error("RouteMessage timeout")
	}
}

func (h *Host) HasRoom(room string) bool {
	h.roomsMu.RLock()
	defer h.roomsMu.RUnlock()
	_, ok := h.Rooms[room]
	return ok
}

func (h *Host) CreateRoom(room, password string) error {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()
	_, ok := h.Rooms[room]
	if ok {
		return fmt.Errorf("room %s already exist", room)
	}
	h.Rooms[room] = &database.RoomEntry{
		Password:  password,
		Timestamp: timestamppb.New(time.Now()),
	}
	return h.Database.UpdateRooms(h.Rooms)
}

func (h *Host) CheckRoomPassword(room, password string) error {
	h.roomsMu.RLock()
	defer h.roomsMu.RUnlock()
	_, ok := h.Rooms[room]
	if !ok {
		return fmt.Errorf("room %s does not exist", room)
	}
	if h.Rooms[room].Password != password {
		return fmt.Errorf("invalid password")
	}
	return nil
}

func (h *Host) resetRoomForConnectedUsers(room string) {
	h.usersMu.RLock()
	defer h.usersMu.RUnlock()
	for _, u := range h.users {
		if u.CurrentRoom != room {
			continue
		}

		if err := u.UpdateCurrentRoom("default"); err != nil {
			h.Log.Error(err)
			continue
		}
		_ = u.WriteLine(aurora.Sprintf("you are now in room %s.", aurora.Blue("default")))
	}
}

func (h *Host) cleanupRooms() error {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()
	pastMarker := time.Now().Add(-3 * time.Hour)
	update := false
	for roomName, roomEntry := range h.Rooms {
		if roomEntry.Timestamp != nil && roomEntry.Timestamp.AsTime().Before(pastMarker) {
			h.Log.Infof("removing room %s", roomName)
			h.resetRoomForConnectedUsers(roomName)
			delete(h.Rooms, roomName)
			update = true
		}
	}
	if !update {
		return nil
	}
	return h.Database.UpdateRooms(h.Rooms)
}

func (h *Host) JoinRoomAnnouncement(u *User) {
	if h.DisableUserAnnouncements {
		return
	}
	h.RoomAnnouncement(u.CurrentRoom, aurora.Sprintf("%s joined the room.", u.RenderName(false)))
}

func (h *Host) LeftRoomAnnouncement(u *User) {
	if h.DisableUserAnnouncements {
		return
	}
	h.RoomAnnouncement(u.CurrentRoom, aurora.Sprintf("%s left the room.", u.RenderName(false)))
}

func (h *Host) ServerInfo() string {
	h.usersMu.RLock()
	userCount := len(h.users)
	h.usersMu.RUnlock()

	verb := "are"
	usersString := "users"
	if userCount == 1 {
		verb = "is"
		usersString = "user"
	}
	return aurora.Sprintf("there %s currently %d %s online.", verb, aurora.Cyan(userCount), usersString)
}

func (h *Host) ListUsersForUser(from *User) error {
	dbUsers, err := h.Database.DumpUsers()
	if err != nil {
		return err
	}
	h.usersMu.RLock()
	defer h.usersMu.RUnlock()
	for _, u := range dbUsers {
		name := aurora.Cyan(u.Name)
		if from.Name == u.Name {
			name = aurora.Magenta(u.Name)
		}
		star := aurora.Yellow("*")
		if _, online := h.users[u.Name]; online {
			star = aurora.Green("*")
		}
		renderStr := aurora.Sprintf("%s %s (%s)", star, name, aurora.Blue(u.CurrentRoom))
		if err := from.WriteLine(renderStr); err != nil {
			return err
		}
	}
	return nil
}

func (h *Host) ShowRecentMessages(u *User, skipAnnouncements bool) error {
	oldMessages, err := h.Database.GetRecentMessagesForUserAndRoom(u.Id, u.CurrentRoom)
	if err != nil {
		return err
	}
	for _, oldMsg := range oldMessages {
		// skip dm history
		if oldMsg.Type == database.MessageType_DIRECT {
			continue
		}
		if skipAnnouncements && oldMsg.Type == database.MessageType_ANNOUNCEMENT {
			continue
		}
		conMsg, err := h.ConvertMessageEntryToMessage(oldMsg)
		if err != nil {
			h.Log.Error(err)
			continue
		}

		if err := u.WriteMessage(conMsg); err != nil {
			h.Log.Error(err)
		}
	}
	return nil
}

func (h *Host) ListRoomsForUser(from *User) error {
	h.roomsMu.RLock()
	defer h.roomsMu.RUnlock()
	for room, roomEntry := range h.Rooms {
		roomEmoji := ":speaking_head:"
		if roomEntry.Password != "" {
			roomEmoji = ":closed_lock_with_key:"
		}
		bulletPoint := aurora.Yellow("*")
		if room == from.CurrentRoom {
			bulletPoint = aurora.Green("*")
		}
		roomInfo := emoji.Sprint(aurora.Sprintf("%s %s %s", bulletPoint, aurora.Blue(room), roomEmoji))
		if err := from.WriteLine(roomInfo); err != nil {
			return err
		}
	}
	return nil
}

func (h *Host) Cleanup() {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()
	h.Log.Info("cleanup started...")
	for {
		select {
		case <-ticker.C:
			err := h.cleanupRooms()
			if err != nil {
				h.Log.Error(err)
			}
		case <-h.cleanupTickerStop:
			h.Log.Info("cleanup stopped")
			return
		}
	}
}

func (h *Host) StopCleanup() {
	h.cleanupTickerStop <- struct{}{}
}

func NewHost(log *logrus.Logger, db *database.Database, rooms map[string]*database.RoomEntry) *Host {
	return &Host{
		Log:                        log,
		users:                      make(map[string]*User),
		msgChan:                    make(chan Message, 50),
		Database:                   db,
		Rooms:                      rooms,
		DisableUserAnnouncements:   true,
		DisableRoomAnnouncements:   false,
		DisableServerAnnouncements: true,
		cleanupTickerStop:          make(chan struct{}, 1),
	}
}
