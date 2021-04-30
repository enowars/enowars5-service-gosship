package chat

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

var Commands = []*Command{
	{
		Prefix:  "help",
		Aliases: []string{"h"},
		Help:    "",
		Handler: func(h *Host, msg *CommandMessage) error {
			help := "this is the help command, yolo!"
			return msg.From.WriteLine(help)
		},
	},
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
