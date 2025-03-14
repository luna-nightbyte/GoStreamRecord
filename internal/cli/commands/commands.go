package commands

import (
	"GoStreamRecord/internal/cli/color"
	"fmt"
)

var passwordWasReset bool

var BinaryName = "GoStreamRecord"

const (
	STARTUP_RESET_PWD       = "reset-pwd"
	STARTUP_ADD_USER        = "add-user"
	STARTUP_GEN_COOKIE_KEY  = "gen-cookie-token"
	STARTUP_GEN_SESSION_KEY = "gen-session-token"
	STARTUP_HELP            = "help"
)

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
	Commands[STARTUP_RESET_PWD] = Command{
		Name:    STARTUP_RESET_PWD,
		Args:    "<username> <new-password>",
		Execute: ResetPwd,
	}
	Commands[STARTUP_ADD_USER] = Command{
		Name:    STARTUP_ADD_USER,
		Args:    "<username> <new-password>",
		Execute: addNewUser,
	}
	Commands[STARTUP_GEN_COOKIE_KEY] = Command{
		Name:    STARTUP_GEN_COOKIE_KEY,
		Args:    "<lenght(int)>",
		Execute: generateCookieToken,
	}
	Commands[STARTUP_GEN_SESSION_KEY] = Command{
		Name:    STARTUP_GEN_SESSION_KEY,
		Args:    "<lenght(int)>",
		Execute: generateSessionKey,
	}
	Commands[STARTUP_HELP] = Command{
		Name:    STARTUP_HELP,
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
