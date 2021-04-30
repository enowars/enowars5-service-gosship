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
	To   string
}

func (d *DirectMessage) String() string {
	return fmt.Sprintf("DirectMessage[to=%s][from=%s]%s", d.To, d.From.Name, d.rawMessage.String())
}

func (d *DirectMessage) RenderFor(u *User) string {
	if u.Name == d.From.Name {
		return aurora.Sprintf("%s[%s]: %s", aurora.Yellow("**"), aurora.Magenta(d.From.Name), d.Body)
	}
	return aurora.Sprintf("%s[%s]: %s", aurora.Yellow("**"), d.From.RenderName(), d.Body)
}

func ParseDirectMessage(args []string, from *User) (Message, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("invalid direct message command")
	}
	return &DirectMessage{
		rawMessage: newRawMessage(strings.Join(args[1:], " ")),
		From:       from,
		To:         args[0],
	}, nil
}

func ParseMessage(m string, from *User) (Message, error) {
	if strings.HasPrefix(m, "/") {
		args := strings.Fields(m)
		switch strings.ToLower(args[0]) {
		case "/dm":
			return ParseDirectMessage(args[1:], from)
		}
		return nil, fmt.Errorf("command not found")
	}
	return &PublicMessage{
		rawMessage: newRawMessage(m),
		From:       from,
		Room:       from.CurrentRoom,
	}, nil
}
