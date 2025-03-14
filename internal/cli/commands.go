package cli

import (
	"GoStreamRecord/internal/cli/color"
	"fmt"
)

var passwordWasReset bool
var BinaryName = "GoStreamRecord"
// Command represents a CLI command with its name, arguments, and execution function.
type Command struct {
	Name    string
	Args    string
	Execute func(args []string) (string, error)
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
		Name:    "add-user",
		Args:    "<username> <new-password>",
		Execute: addNewUser,
	}
	Commands["gen-cookie-token"] = Command{
		Name:    "gen-cookie-token",
		Args:    "<lenght(int)>",
		Execute: generateCookieToken,
	}
	Commands["gen-session-token"] = Command{
		Name:    "gen-session-token",
		Args:    "<lenght(int)>",
		Execute: generateSessionKey,
	}
	Commands["help"] = Command{
		Name:    "help",
		Args:    "Shows this menu",
		Execute: PrintUsage,
	}
}

// args are not used here. input "nil"
func PrintUsage(args []string) (string, error) {

	// TODO: cli.PrintUsage()
	fmt.Println()
	color.Println("Bgrey", "Usage:")
	for _, cmd := range Commands {
		color.Print("cyan", " - ./"+BinaryName)
		color.Print("white", " "+cmd.Name)
		color.Println("Bwhite", " "+cmd.Args)
	}

	color.Println("Bgrey", "\nOtherwise run the server without any arguments.")

	return "", nil
}
