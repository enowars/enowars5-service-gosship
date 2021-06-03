package chat

import (
	"errors"
	"gosship/pkg/database"
	"gosship/pkg/terminal"
	"io"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/google/uuid"
	"github.com/logrusorgru/aurora/v3"
	gossh "golang.org/x/crypto/ssh"
)

var ErrFingerprintDoesNotMatch = errors.New("the public key does not match with the username")
var ErrFingerprintAlreadyRegistered = errors.New("the public key is already used")
var ErrDummyUser = errors.New("user is a dummy")

type User struct {
	Id              string
	Session         ssh.Session
	Name            string
	Term            *terminal.Terminal
	CurrentRoom     string
	Fingerprint     string
	db              *database.Database
	Dummy           bool
	LastDmRecipient string
}

func NewUser(db *database.Database, session ssh.Session) (*User, error) {
	name := session.User()
	userId, userEntry, err := db.FindUserByPredicate(func(entry *database.UserEntry) bool {
		return entry.Name == name
	})
	if err != nil {
		return nil, err
	}

	fingerprint := gossh.FingerprintSHA256(session.PublicKey())
	if userId != "" && userEntry.Fingerprint != fingerprint {
		return nil, ErrFingerprintDoesNotMatch
	}

	if userId == "" {
		fingerprintUserId, _, err := db.FindUserByPredicate(func(entry *database.UserEntry) bool {
			return entry.Fingerprint == fingerprint
		})
		if err != nil {
			return nil, err
		}
		if fingerprintUserId != "" {
			return nil, ErrFingerprintAlreadyRegistered
		}
	}

	u := &User{
		Session:     session,
		Name:        name,
		Term:        terminal.New(session),
		CurrentRoom: "default",
		Fingerprint: fingerprint,
		db:          db,
		Id:          "",
		Dummy:       false,
	}

	u.UpdatePrompt()

	if userId != "" {
		u.Id = userId
		u.CurrentRoom = userEntry.CurrentRoom
	} else {
		u.Id = uuid.NewString()
	}

	err = u.DBUpdate()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) UpdatePrompt() {
	u.Term.SetPrompt(aurora.Sprintf("[%s]: ", aurora.Magenta(u.Name)))
}

func (u *User) WriteLine(line string) error {
	if u.Dummy {
		return ErrDummyUser
	}
	_, err := io.WriteString(u.Term, line+"\n")
	if err == io.EOF {
		return nil
	}
	return err
}

func (u *User) WriteMessage(msg Message) error {
	if u.Dummy {
		return ErrDummyUser
	}
	return u.WriteLine(msg.RenderFor(u))
}

func (u *User) RenderName(self bool) string {
	if self {
		return aurora.Magenta(u.Name).String()
	}
	return aurora.Cyan(u.Name).String()
}

func (u *User) DBUpdate() error {
	return u.db.AddOrUpdateUser(u.Id, &database.UserEntry{
		Fingerprint: u.Fingerprint,
		Name:        u.Name,
		CurrentRoom: u.CurrentRoom,
	})
}

func (u *User) UpdateCurrentRoom(room string) error {
	u.CurrentRoom = strings.ToLower(room)
	return u.DBUpdate()
}
