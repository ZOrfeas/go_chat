package cli

import (
	"fmt"
	"strings"
)

type Callback func(args []string) error

type Command struct {
	Token   string
	Details string
	Action  Callback
}

type Cli struct {
	name     string
	details  string
	commands map[string]*Command
}

func NewCli(name, details string) *Cli {
	tmpCli := &Cli{
		name:     name,
		details:  details,
		commands: make(map[string]*Command),
	}
	tmpCli.AddCommand("help", "Displays this message", tmpCli.Help)
	return tmpCli
}

func (cli *Cli) exists(token string) bool {
	_, ok := cli.commands[token]
	return ok
}

func (cli *Cli) AddCommand(token, details string, action Callback) error {
	if cli.exists(token) {
		return fmt.Errorf("Command token '%s' declared twice", token)
	}

	newCommand := &Command{
		Token:   token,
		Details: details,
		Action:  action,
	}
	cli.commands[token] = newCommand
	return nil
}
func (cli *Cli) Help(args []string) error {
	var helpBuilder strings.Builder
	helpBuilder.WriteString(cli.name)
	helpBuilder.WriteString(" - ")
	helpBuilder.WriteString(cli.details)
	helpBuilder.WriteString("\n")
	helpBuilder.WriteString("\tUsage: ")
	helpBuilder.WriteString(cli.name)
	helpBuilder.WriteString(" [command] [arguments] [command] [arguments] ...\n")
	helpBuilder.WriteString("Commands:\n")
	for cmdName, cmd := range cli.commands {
		helpBuilder.WriteString("\t" + cmdName)
		helpBuilder.WriteString("\t\t" + cmd.Details + "\n")
	}
	fmt.Print(helpBuilder.String())
	return nil
}

func (cli *Cli) Run(args []string, dflt Callback) error {
	atLeastOneRun := false
	for i := 0; i < len(args); i += 1 {
		if cli.exists(args[i]) {
			// captures the command arguments
			subArgs := []string{}
			j := i + 1
			for ; (j < len(args)) && !cli.exists(args[j]); j += 1 {
				subArgs = append(subArgs, args[j])
			}

			// call to the callback and gives it its arguments
			err := cli.commands[args[i]].Action(subArgs)
			atLeastOneRun = true
			if err != nil {
				return err
			}
			i = j - 1
		}
	}
	if !atLeastOneRun {
		dflt([]string{})
	}
	return nil
}
