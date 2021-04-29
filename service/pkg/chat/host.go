package chat

import (
	"fmt"
	"io"
	"sync"

	"github.com/gliderlabs/ssh"
	"github.com/logrusorgru/aurora/v3"
	"github.com/sirupsen/logrus"
	gossh "golang.org/x/crypto/ssh"
)

const title = "              _____ _____ _    _ _       \n             / ____/ ____| |  | (_)      \n   __ _  ___| (___| (___ | |__| |_ _ __  \n  / _` |/ _ \\\\___ \\\\___ \\|  __  | | '_ \\ \n | (_| | (_) |___) |___) | |  | | | |_) |\n  \\__, |\\___/_____/_____/|_|  |_|_| .__/ \n   __/ |                          | |    \n  |___/                           |_|    "

type Host struct {
	log     *logrus.Logger
	mu      sync.Mutex
	users   map[string]*User
	msgChan chan Message
}

func (h *Host) HandleNewSession(session ssh.Session) {
	code := h.handleNewSessionWithExitCode(session)
	if err := session.Exit(code); err != nil {
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
	return true
}

func (h *Host) handleNewSessionWithExitCode(session ssh.Session) int {
	_, _, isPty := session.Pty()
	if !isPty {
		_, err := io.WriteString(session, "No PTY requested.\n")
		if err != nil {
			h.log.Error(err)
			return -1
		}
	}
	h.log.Printf("new session: user=%s [%s]\n", session.User(), gossh.FingerprintLegacyMD5(session.PublicKey()))
	u := NewUser(session)
	if !h.AddUser(u) {
		_, err := fmt.Fprintf(u.Term, "%s is already logged in!\n\n", aurora.Red(u.Name))
		if err != nil {
			h.log.Error(err)
			return -1
		}
		return -1
	}
	_, err := fmt.Fprintf(u.Term, "%s\n\n🦄 Welcome %s!\n\n", aurora.Green(title), aurora.Magenta(u.Name))
	if err != nil {
		h.log.Error(err)
		return -1
	}
	_, _ = u.Term.ReadLine()
	//TODO: Parse Message
	//TODO: send to h.msgChan
	return 0
}

func (h *Host) HandlePublicKey(ctx ssh.Context, key ssh.PublicKey) bool {
	h.log.Printf("new connection (%s) with key type: %s\n", ctx.RemoteAddr(), key.Type())
	return true
}

func (h *Host) Serve() {
	for msg := range h.msgChan {
		h.log.Println(msg)
	}
}

func NewHost(log *logrus.Logger) *Host {
	return &Host{
		log:     log,
		users:   make(map[string]*User),
		msgChan: make(chan Message),
	}
}
