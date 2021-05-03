package chat

import (
	"fmt"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
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
	return fmt.Sprintf("[%s]: %s", m.Timestamp.Format(time.RFC3339), stripansi.Strip(m.Body))
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
	return aurora.Sprintf(aurora.Yellow("-> %s"), a.Body)
}

type RoomAnnouncementMessage struct {
	*rawMessage
	Room string
}

func NewRoomAnnouncementMessage(room, body string) *RoomAnnouncementMessage {
	return &RoomAnnouncementMessage{
		rawMessage: newRawMessage(body),
		Room:       room,
	}
}

func (r *RoomAnnouncementMessage) String() string {
	return fmt.Sprintf("RoomAnnouncementMessage[room=%s]%s", r.Room, r.rawMessage.String())
}

func (r *RoomAnnouncementMessage) RenderFor(u *User) string {
	return aurora.Sprintf(aurora.Yellow("-> %s"), r.Body)
}

type PublicMessage struct {
	*rawMessage
	From *User
	Room string
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
	From       *User
	To         string
	ToResolved *User
}

func (d *DirectMessage) String() string {
	return fmt.Sprintf("DirectMessage[to=%s][from=%s]%s", d.To, d.From.Name, d.rawMessage.String())
}

func (d *DirectMessage) RenderFor(u *User) string {
	if u.Name == d.From.Name {
		return aurora.Sprintf("%s [%s]: %s", aurora.Yellow("*private*"), aurora.Magenta(d.From.Name), d.Body)
	}
	return aurora.Sprintf("%s [%s]: %s", aurora.Yellow("*private*"), d.From.RenderName(), d.Body)
}

type CommandMessage struct {
	*rawMessage
	From *User
	Cmd  string
	Args []string
}

func (c *CommandMessage) String() string {
	return fmt.Sprintf("CommandMessage[from=%s][cmd=%s][args=%s]%s", c.From.Name, c.Cmd, c.Args, c.rawMessage.String())
}

func (c *CommandMessage) RenderFor(u *User) string {
	return c.String()
}

func NewCommandMessage(rawBody, cmd string, args []string, from *User) *CommandMessage {
	return &CommandMessage{
		rawMessage: newRawMessage(rawBody),
		From:       from,
		Cmd:        cmd,
		Args:       args,
	}
}

func ParseDirectMessage(args []string, from *User) (Message, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("invalid direct message command")
	}
	if args[0] == from.Name {
		return nil, fmt.Errorf("you cannot send a direct message to yourself")
	}
	return &DirectMessage{
		rawMessage: newRawMessage(strings.Join(args[1:], " ")),
		From:       from,
		To:         args[0],
		ToResolved: nil,
	}, nil
}

func ParseMessage(m string, from *User) (Message, error) {
	if strings.HasPrefix(m, "/") {
		args := strings.Fields(m)
		cmd := strings.ToLower(args[0])[1:]
		if len(cmd) == 0 {
			return nil, fmt.Errorf("invalid command")
		}
		if cmd == "dm" {
			return ParseDirectMessage(args[1:], from)
		}
		return NewCommandMessage(m, cmd, args[1:], from), nil
	}
	return &PublicMessage{
		rawMessage: newRawMessage(m),
		From:       from,
		Room:       from.CurrentRoom,
	}, nil
}
