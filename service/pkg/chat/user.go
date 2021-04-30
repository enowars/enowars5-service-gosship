package chat

import (
	"gosship/pkg/terminal"
	"io"

	"github.com/gliderlabs/ssh"
	"github.com/logrusorgru/aurora/v3"
	gossh "golang.org/x/crypto/ssh"
)

type User struct {
	Session     ssh.Session
	Name        string
	Term        *terminal.Terminal
	CurrentRoom Room
	Fingerprint string
}

func NewUser(session ssh.Session) *User {
	name := session.User()
	prompt := aurora.Sprintf("[%s]: ", aurora.Magenta(name))
	u := &User{
		Session:     session,
		Name:        name,
		Term:        terminal.New(session, prompt),
		CurrentRoom: "default",
		Fingerprint: gossh.FingerprintLegacyMD5(session.PublicKey()),
	}
	return u
}

func (u *User) WriteLine(line string) error {
	_, err := io.WriteString(u.Term, line+"\n")
	if err == io.EOF {
		return nil
	}
	return err

}

func (u *User) WriteMessage(msg Message) error {
	return u.WriteLine(msg.RenderFor(u))
}

func (u *User) RenderName() string {
	return aurora.Cyan(u.Name).String()
}
