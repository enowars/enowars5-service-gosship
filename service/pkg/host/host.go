package host

import (
	"io"
	"log"

	gossh "golang.org/x/crypto/ssh"

	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
)

type Host struct {
	log *logrus.Logger
}

func (h *Host) HandleNewSession(session ssh.Session) {
	_, _, isPty := session.Pty()
	if !isPty {
		_, err := io.WriteString(session, "No PTY requested.\n")
		if err != nil {
			h.log.Error(err)
			if err := session.Exit(1); err != nil {
				h.log.Error(err)
			}
			return
		}
	}
	log.Printf("new session: user=%s [%s]\n", session.User(), gossh.FingerprintLegacyMD5(session.PublicKey()))
	_, err := io.WriteString(session, "hello!\n")
	if err != nil {
		return
	}
	err = session.Exit(0)
	if err != nil {
		h.log.Error(err)
	}
}

func (h *Host) HandlePublicKey(ctx ssh.Context, key ssh.PublicKey) bool {
	h.log.Printf("new connection (%s) with key type: %s\n", ctx.RemoteAddr(), key.Type())
	return true
}

func New(log *logrus.Logger) *Host {
	return &Host{log: log}
}
