package chat

import (
	"fmt"
	"gosship/pkg/database"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/logrusorgru/aurora/v3"
)

type Command struct {
	Prefix  string
	Args    string
	Aliases []string
	Help    string
	Handler func(h *Host, msg *CommandMessage) error
}

func (c *Command) MatchCommand(cmdPrefix string) bool {
	if c.Prefix == cmdPrefix {
		return true
	}
	for _, a := range c.Aliases {
		if a == cmdPrefix {
			return true
		}
	}
	return false
}

func (c *Command) GetTableRow() table.Row {
	aliasesTxt := ""
	for _, a := range c.Aliases {
		aliasesTxt += "/" + a + ", "
	}
	if aliasesTxt != "" {
		aliasesTxt = aliasesTxt[0 : len(aliasesTxt)-2]
	}
	cmd := "/" + c.Prefix
	if c.Args != "" {
		cmd += " " + c.Args
	}
	return []interface{}{cmd, aliasesTxt, c.Help}
}

var helpTableRender string

var Commands = []*Command{
	{
		Prefix: "dm",
		Args:   "[user] [msg]",
		Help:   "send a direct message to a user",
		Handler: func(h *Host, msg *CommandMessage) error {
			dm, err := ParseDirectMessage(msg.Args, msg.From)
			if err != nil {
				return err
			}
			h.RouteMessage(dm)
			return nil
		},
	},
	{
		Prefix:  "help",
		Aliases: []string{"h", "?"},
		Help:    "show the help for all available commands",
		Handler: func(h *Host, msg *CommandMessage) error {
			return msg.From.WriteLine(helpTableRender)
		},
	},
	{
		Prefix:  "exit",
		Aliases: []string{"quit", "q"},
		Help:    "leave the chat",
		Handler: func(h *Host, msg *CommandMessage) error {
			_ = msg.From.WriteLine("bye!")
			return msg.From.Session.Close()
		},
	},
	{
		Prefix:  "info",
		Aliases: []string{"i"},
		Help:    "info about the logged-in user",
		Handler: func(h *Host, msg *CommandMessage) error {
			info := table.NewWriter()
			info.AppendRow(table.Row{"Database ID", msg.From.Id})
			info.AppendRow(table.Row{"Name", msg.From.Name})
			info.AppendRow(table.Row{"Fingerprint", msg.From.Fingerprint})
			info.AppendRow(table.Row{"Current Room", msg.From.CurrentRoom})
			return msg.From.WriteLine(info.Render())
		},
	},
	{
		Prefix:  "reply",
		Args:    "[msg]",
		Aliases: []string{"r"},
		Help:    "reply to your last direct message",
		Handler: func(h *Host, msg *CommandMessage) error {
			if msg.From.LastDmRecipient == "" {
				return fmt.Errorf("no current message conversation")
			}
			dm, err := ParseDirectMessage(append([]string{msg.From.LastDmRecipient}, msg.Args...), msg.From)
			if err != nil {
				return err
			}
			h.RouteMessage(dm)
			return nil
		},
	},
	{
		Prefix: "history",
		Args:   "[user]",
		Help:   "show the direct message history",
		Handler: func(h *Host, msg *CommandMessage) error {
			if len(msg.Args) == 0 {
				return fmt.Errorf("user argument is missing")
			}

			uid, _, err := h.Database.FindUserByPredicate(func(entry *database.UserEntry) bool {
				return entry.Name == msg.Args[0]
			})
			if err != nil {
				h.Log.Error(err)
				return fmt.Errorf("user not found")
			}
			if uid == "" {
				return fmt.Errorf("user not found")
			}

			dms, err := h.Database.GetRecentDirectMessagesForUser(msg.From.Id, uid)
			if err != nil {
				return fmt.Errorf("messages not found")
			}
			if len(dms) == 0 {
				return fmt.Errorf("no messages found")
			}
			for _, rawDm := range dms {
				dm, err := h.ConvertMessageEntryToMessage(rawDm)
				if err != nil {
					h.Log.Error(err)
					continue
				}
				err = msg.From.WriteMessage(dm)
				if err != nil {
					h.Log.Error(err)
					return err
				}
			}
			return nil
		},
	},
	{
		Prefix: "shrug",
		Help:   "¯\\_(ツ)_/¯",
		Handler: func(h *Host, msg *CommandMessage) error {
			h.RouteMessage(NewPublicMessage("¯\\_(ツ)_/¯", msg.From))
			return nil
		},
	},
	{
		Prefix: "rename",
		Args:   "[new name]",
		Help:   "change your username",
		Handler: func(h *Host, msg *CommandMessage) error {
			if len(msg.Args) == 0 {
				return fmt.Errorf("new name argument is missing")
			}
			newName := msg.Args[0]
			if newName == msg.From.Name {
				return fmt.Errorf("the new name is identical with your current name")
			}
			err := h.Database.RenameUser(msg.From.Id, newName)
			if err != nil {
				return err
			}
			h.RoomAnnouncement(msg.From.CurrentRoom, aurora.Sprintf("%s changed their name to %s", msg.From.RenderName(false), aurora.Cyan(newName)))
			msg.From.Name = newName
			msg.From.UpdatePrompt()
			return nil
		},
	}, {
		Prefix: "create",
		Args:   "[room] <password>",
		Help:   "create a new room",
		Handler: func(h *Host, msg *CommandMessage) error {
			if len(msg.Args) == 0 {
				return fmt.Errorf("room argument is missing")
			}
			roomName := msg.Args[0]
			password := ""
			if len(msg.Args) == 2 {
				password = msg.Args[1]
			}

			if err := h.CreateRoom(roomName, password); err != nil {
				return err
			}

			h.Announcement(aurora.Sprintf("room %s was created by %s.", aurora.Blue(roomName), msg.From.RenderName(false)))
			h.LeftRoomAnnouncement(msg.From)
			if err := msg.From.UpdateCurrentRoom(roomName); err != nil {
				return err
			}
			return msg.From.WriteLine(aurora.Sprintf("you are now in room %s.", aurora.Blue(roomName)))
		},
	},
	{
		Prefix:  "join",
		Args:    "<room> <password>",
		Aliases: []string{"j"},
		Help:    "join a room",
		Handler: func(h *Host, msg *CommandMessage) error {
			roomName := "default"
			if len(msg.Args) > 0 {
				roomName = msg.Args[0]
			}
			password := ""
			if len(msg.Args) == 2 {
				password = msg.Args[1]
			}
			if msg.From.CurrentRoom == roomName {
				return fmt.Errorf("you are already in room %s", roomName)
			}

			if err := h.CheckRoomPassword(roomName, password); err != nil {
				return err
			}

			h.LeftRoomAnnouncement(msg.From)
			if err := msg.From.UpdateCurrentRoom(roomName); err != nil {
				return err
			}
			_ = msg.From.WriteLine(aurora.Sprintf("you are now in room %s.", aurora.Blue(roomName)))
			if err := h.ShowRecentMessages(msg.From, true); err != nil {
				return err
			}
			h.JoinRoomAnnouncement(msg.From)
			return nil
		},
	},
	{
		Prefix: "users",
		Help:   "list users on the server",
		Handler: func(h *Host, msg *CommandMessage) error {
			if err := msg.From.WriteLine(h.ServerInfo()); err != nil {
				return err
			}
			return h.ListUsersForUser(msg.From)
		},
	},
	{
		Prefix: "rooms",
		Help:   "list rooms on the server",
		Handler: func(h *Host, msg *CommandMessage) error {
			if err := msg.From.WriteLine("the following rooms exist on the server:"); err != nil {
				return err
			}
			return h.ListRoomsForUser(msg.From)
		},
	},
}

func init() {
	helpTable := table.NewWriter()
	helpTable.AppendHeader(table.Row{"Command", "Aliases", "Help"})
	helpTable.AppendSeparator()
	for _, cmd := range Commands {
		helpTable.AppendRow(cmd.GetTableRow())
		helpTable.AppendSeparator()
	}
	helpTableRender = helpTable.Render()
}

func FindCommand(cmdPrefix string) *Command {
	for _, cmd := range Commands {
		if !cmd.MatchCommand(cmdPrefix) {
			continue
		}
		return cmd
	}
	return nil
}
