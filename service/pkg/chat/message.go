package chat

import (
	"fmt"
	"strings"
	"time"

	"github.com/logrusorgru/aurora/v3"
)

type Message interface {
	fmt.Stringer
	RenderFor(u *User) string
}

type rawMessage struct {
	Timestamp time.Time
	Body      string
}

func newRawMessage(body string) *rawMessage {
	return &rawMessage{
		Body:      body,
		Timestamp: time.Now(),
	}
}

func (m *rawMessage) String() string {
	return fmt.Sprintf("[%s]: %s", m.Timestamp.Format(time.RFC3339), m.Body)
}

func (m *rawMessage) RenderFor(u *User) string {
	return m.String()
}

type AnnouncementMessage struct {
	*rawMessage
}

func NewAnnouncementMessage(body string) *AnnouncementMessage {
	return &AnnouncementMessage{newRawMessage(body)}
}

func (a *AnnouncementMessage) String() string {
	return fmt.Sprintf("AnnounceMessage%s", a.rawMessage.String())
}

func (a *AnnouncementMessage) RenderFor(u *User) string {
	return aurora.Sprintf("%s", aurora.Gray(12, a.Body))
}

type PublicMessage struct {
	*rawMessage
	From *User
	Room Room
}

func (p *PublicMessage) String() string {
	return fmt.Sprintf("PublicMessage[room=%s][from=%s]%s", p.Room, p.From.Name, p.rawMessage.String())
}

func (p *PublicMessage) RenderFor(u *User) string {
	if u.Name == p.From.Name {
		return aurora.Sprintf("[%s]: %s", aurora.Magenta(p.From.Name), p.Body)
	}
	return aurora.Sprintf("[%s]: %s", p.From.RenderName(), p.Body)
}

type DirectMessage struct {
	*rawMessage
	From *User
	To   *User
}

func ParseMessage(m string, from *User) (Message, error) {
	if strings.HasPrefix(m, "/") {
		return nil, fmt.Errorf("commands not implemented yet")
	}
	return &PublicMessage{
		rawMessage: newRawMessage(m),
		From:       from,
		Room:       from.CurrentRoom,
	}, nil
}
