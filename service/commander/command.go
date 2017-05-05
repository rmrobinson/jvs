package commander

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
)

type Command interface {
	Name() string
	Help() string

	Execute(ctx context.Context, args []string) (string, error)
}

type RootCommand struct {
	subCommands map[string]Command
}

func (c *RootCommand) AddSubCommand(cmd Command) {
	if c.subCommands == nil {
		c.subCommands = make(map[string]Command)
	}
	c.subCommands[cmd.Name()] = cmd
}

func (c *RootCommand) RemoveSubCommand(name string) {
	delete(c.subCommands, name)
}

func (c *RootCommand) Execute(ctx context.Context, args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("Invalid command specified")
	}

	if args[0] == "help" || args[0] == "?" {
		var ret string

		for _, subCommand := range c.subCommands {
			ret += fmt.Sprintf("%s\t%s\n", subCommand.Name(), subCommand.Help())
		}

		return ret, nil
	}

	if child, ok := c.subCommands[args[0]]; ok {
		return child.Execute(ctx, args[1:])
	}

	return "", errors.New("Unknown command specified")
}
