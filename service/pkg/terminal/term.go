package terminal

import (
	"github.com/gliderlabs/ssh"

	term "github.com/shazow/ssh-chat/sshd/terminal"
)

type Terminal struct {
	*term.Terminal
	Session ssh.Session
}

func (t *Terminal) handleWinCh() {
	ptyReq, winCh, _ := t.Session.Pty()
	_ = t.SetSize(ptyReq.Window.Width, ptyReq.Window.Height)
	go func() {
		for win := range winCh {
			_ = t.SetSize(win.Width, win.Height)
		}
	}()
}

func New(session ssh.Session, prompt string) *Terminal {
	t := &Terminal{
		Terminal: term.NewTerminal(session, prompt),
		Session:  session,
	}
	t.SetEnterClear(true)
	go t.handleWinCh()
	return t
}
