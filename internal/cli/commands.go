package cli

import (
	"GoStreamRecord/internal/cli/color"
	"fmt"
)

var passwordWasReset bool

// Command represents a CLI command with its name, arguments, and execution function.
type Command struct {
	Name    string
	Args    string
	Execute func(args []string)
}

// Global command registry.
var Commands = map[string]Command{}

func init() {
	// Register available commands.
	Commands["reset-pwd"] = Command{
		Name:    "reset-pwd",
		Args:    "<username> <new-password>",
		Execute: resetPwdCommand,
	}
	Commands["add-user"] = Command{
		Name:    "reset-pwd",
		Args:    "<username> <new-password>",
		Execute: addNewUser,
	}
}

func PrintUsage() {

	// TODO: cli.PrintUsage()
	fmt.Println()
	color.Println("Bgrey", "Usage:")
	for _, cmd := range Commands {
		color.Print("cyan", " - ./GoStreamRecord")
		color.Print("white", " "+cmd.Name)
		color.Println("Bwhite", " "+cmd.Args)
	}

	color.Println("Bgrey", "\nOtherwise run the server without any arguments.")

}
