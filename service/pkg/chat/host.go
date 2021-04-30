package chat

import (
	"fmt"
	"io"
	"sync"

	"github.com/gliderlabs/ssh"
	"github.com/logrusorgru/aurora/v3"
	"github.com/sirupsen/logrus"
)

const title = "              _____ _____ _    _ _       \n             / ____/ ____| |  | (_)      \n   __ _  ___| (___| (___ | |__| |_ _ __  \n  / _` |/ _ \\\\___ \\\\___ \\|  __  | | '_ \\ \n | (_| | (_) |___) |___) | |  | | | |_) |\n  \\__, |\\___/_____/_____/|_|  |_|_| .__/ \n   __/ |                          | |    \n  |___/                           |_|    "

type Host struct {
	log     *logrus.Logger
	mu      sync.RWMutex
	users   map[string]*User
	msgChan chan Message
}

func (h *Host) HandleNewSession(session ssh.Session) {
	code := 0
	err := h.handleNewSessionWithError(session)
	if err != nil {
		h.log.Error(err)
		code = -1
	}
	if err := session.Exit(code); err != nil && err != io.EOF {
		h.log.Error(err)
	}
}

func (h *Host) AddUser(u *User) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.users[u.Name]; ok {
		return false
	}
	h.users[u.Name] = u
	h.log.Printf("[%s] added\n", u.Name)
	return true
}

func (h *Host) RemoveUser(name string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.users, name)
	h.log.Printf("[%s] removed\n", name)
}

func (h *Host) handleNewSessionWithError(session ssh.Session) error {
	_, _, isPty := session.Pty()
	if !isPty {
		_, err := io.WriteString(session, "No PTY requested.\n")
		if err != nil {
			return err
		}
		return fmt.Errorf("no PTY requested")
	}
	u := NewUser(session)
	h.log.Printf("[%s] new session: fingerprint=(%s)\n", u.Name, u.Fingerprint)
	if !h.AddUser(u) {
		err := u.WriteLine(aurora.Sprintf("%s is already logged in!\n", aurora.Red(u.Name)))
		if err != nil {
			return err
		}
		return fmt.Errorf("[%s] already logged in", u.Name)
	}
	defer h.RemoveUser(u.Name)
	err := u.WriteLine(aurora.Sprintf("%s\n\n🦄 Welcome %s!\n", aurora.Green(title), aurora.Magenta(u.Name)))
	if err != nil {
		return err
	}
	h.Announcement(aurora.Sprintf("%s joined the room.", u.RenderName()))
	for {
		line, err := u.Term.ReadLine()
		if err != nil {
			if err != io.EOF {
				h.log.Error(err)
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
		h.msgChan <- parsedMessage
	}
	h.Announcement(aurora.Sprintf("%s left the room.", u.RenderName()))
	return nil
}

func (h *Host) HandlePublicKey(ctx ssh.Context, key ssh.PublicKey) bool {
	h.log.Printf("new connection (%s) with key type: %s\n", ctx.RemoteAddr(), key.Type())
	return true
}

func (h *Host) Serve() {
	for msg := range h.msgChan {
		h.log.Println(msg.String())
		switch v := msg.(type) {
		case *PublicMessage:
			h.sendMessageToAllUsersInRoom(v)
		case *AnnouncementMessage:
			h.sendMessageToAllUsers(v)
		case *DirectMessage:
			h.sendMessageToUser(v)
		case *CommandMessage:
			h.handleUserCommand(v)
		default:
			h.log.Error("unknown message type")
		}
	}
}

func (h *Host) sendMessageToAllUsersInRoom(msg *PublicMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, u := range h.users {
		if u.CurrentRoom != msg.Room {
			continue
		}
		err := u.WriteMessage(msg)
		if err != nil {
			h.log.Error(err)
		}
	}
}

func (h *Host) sendMessageToAllUsers(msg *AnnouncementMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, u := range h.users {
		err := u.WriteMessage(msg)
		if err != nil {
			h.log.Error(err)
		}
	}
}

func (h *Host) sendMessageToUser(msg *DirectMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	to, ok := h.users[msg.To]
	if !ok {
		err := msg.From.WriteLine(aurora.Sprintf(aurora.Yellow("user %s not found on the server."), aurora.Red(msg.To)))
		if err != nil {
			h.log.Error(err)
		}
		return
	}
	for _, u := range []*User{msg.From, to} {
		err := u.WriteMessage(msg)
		if err != nil {
			h.log.Error(err)
		}
	}
}
func (h *Host) handleUserCommand(msg *CommandMessage) {
	cmd := FindCommand(msg.Cmd)
	if cmd == nil {
		err := msg.From.WriteLine(aurora.Sprintf(aurora.Yellow("command %s not found. use /help to list all available commands."), aurora.Red(msg.Cmd)))
		if err != nil {
			h.log.Error(err)
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

func NewHost(log *logrus.Logger) *Host {
	return &Host{
		log:     log,
		users:   make(map[string]*User),
		msgChan: make(chan Message, 10),
	}
}
