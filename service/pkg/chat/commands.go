package chat

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
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
		Prefix: "info",
		Help:   "info about the logged-in user",
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
		Prefix: "rename",
		Help:   "change your username",
		Handler: func(h *Host, msg *CommandMessage) error {
			return msg.From.WriteLine("TODO")
		},
	},
	{
		Prefix: "reply",
		Args:   "[msg]",
		Help:   "reply to your last direct message",
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
			return msg.From.WriteLine("TODO")
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
