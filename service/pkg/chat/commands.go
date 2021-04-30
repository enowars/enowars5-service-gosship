package chat

import (
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Command struct {
	Prefix  string
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
	return []interface{}{c.Prefix, strings.Join(c.Aliases, ", "), c.Help}
}

var helpTableRender string

var Commands = []*Command{
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
		Aliases: []string{"quit"},
		Help:    "leave the chat",
		Handler: func(h *Host, msg *CommandMessage) error {
			_ = msg.From.WriteLine("bye!")
			return msg.From.Session.Close()
		},
	},
}

func init() {
	helpTable := table.NewWriter()
	helpTable.AppendHeader(table.Row{"Command", "Aliases", "Help"})
	helpTable.AppendRow(table.Row{"dm [user] [msg]", "", "send a direct message to a user"})
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
