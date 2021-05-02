package chat

import (
	"gosship/pkg/database"
	"gosship/pkg/terminal"
	"io"

	"github.com/google/uuid"

	"github.com/gliderlabs/ssh"
	"github.com/logrusorgru/aurora/v3"
	gossh "golang.org/x/crypto/ssh"
)

type User struct {
	Id          string
	Session     ssh.Session
	Name        string
	Term        *terminal.Terminal
	CurrentRoom Room
	Fingerprint string
	db          *database.Database
}

func NewUser(db *database.Database, session ssh.Session) (*User, error) {
	fingerprint := gossh.FingerprintLegacyMD5(session.PublicKey())
	userId, userEntry, err := db.FindUserByFingerprint(fingerprint)
	if err != nil {
		return nil, err
	}

	name := session.User()
	prompt := aurora.Sprintf("[%s]: ", aurora.Magenta(name))
	u := &User{
		Session:     session,
		Name:        name,
		Term:        terminal.New(session, prompt),
		CurrentRoom: "default",
		Fingerprint: fingerprint,
		db:          db,
		Id:          "",
	}

	if userId != "" {
		u.Id = userId
		u.CurrentRoom = Room(userEntry.CurrentRoom)
	} else {
		u.Id = uuid.NewString()
	}

	err = db.AddOrUpdateUser(u.Id, &database.UserEntry{
		Fingerprint: u.Fingerprint,
		Name:        u.Name,
		CurrentRoom: string(u.CurrentRoom),
	})
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) WriteLine(line string) error {
	_, err := io.WriteString(u.Term, line+"\n")
	if err == io.EOF {
		return nil
	}
	return err

}

func (u *User) pushDirectMessage(dm *DirectMessage) {
	// TODO
}

func (u *User) WriteMessage(msg Message) error {
	if dm, ok := msg.(*DirectMessage); ok {
		u.pushDirectMessage(dm)
	}
	return u.WriteLine(msg.RenderFor(u))
}

func (u *User) RenderName() string {
	return aurora.Cyan(u.Name).String()
}
