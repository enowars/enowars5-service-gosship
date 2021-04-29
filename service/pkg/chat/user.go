package chat

import (
	"gosship/pkg/terminal"

	"golang.org/x/term"

	"github.com/gliderlabs/ssh"
)

type User struct {
	Session ssh.Session
	Name    string
	Term    *term.Terminal
	//currentRoom *Room
}

func NewUser(session ssh.Session) *User {
	name := session.User()

	u := &User{
		Session: session,
		Name:    name,
		Term:    terminal.New(session),
	}
	u.handleWinCh()
	return u
}

func (u *User) handleWinCh() {
	ptyReq, winCh, _ := u.Session.Pty()
	_ = u.Term.SetSize(ptyReq.Window.Width, ptyReq.Window.Height)
	go func() {
		for win := range winCh {
			_ = u.Term.SetSize(win.Width, win.Height)
		}
	}()
}

//
//func (u *User) handleMessages() {
//	for {
//		line, err := u.terminal.ReadLine()
//		if err == io.EOF {
//			break
//		}
//		checkErr(err, "handleMessages")
//		if u.currentRoom != nil {
//			u.currentRoom.Send(NewMessage(u, line))
//		}
//	}
//	checkErr(u.session.Exit(0))
//}
//
//func (u *User) writeString(msg string) error {
//	_, err := io.WriteString(u.terminal, msg)
//	if err == io.EOF {
//		return nil
//	}
//	return err
//}
//
//func (u *User) SendMessage(msg *Message) {
//	u.writeString(aurora.Sprintf("\r%20s | %s\n", aurora.Green(msg.from.name), msg.msg))
//}

//func (u *User) SetRoom(room *Room) {
//	u.currentRoom = room
//}
