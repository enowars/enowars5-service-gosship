package chat

import (
	"fmt"
	"time"
)

type Message interface {
	isMessage()
	fmt.Stringer
}

type msg struct {
	From      *User
	Body      string
	Timestamp time.Time
}

func (*msg) isMessage() {}
func (m *msg) String() string {
	return m.Body
}

type PublicMessage struct {
	*msg
}

type DirectMessage struct {
	*msg
	To *User
}

func ParseMessage(m string, from *User) Message {
	return &PublicMessage{
		msg: &msg{
			From:      from,
			Body:      m,
			Timestamp: time.Now(),
		},
	}
}
