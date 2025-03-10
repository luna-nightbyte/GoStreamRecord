package startup

var passwordWasReset bool

// Command represents a CLI command with its name, usage, and execution function.
type Command struct {
	Name    string
	Usage   string
	Execute func(args []string)
}

// Global command registry.
var Commands = map[string]Command{}

func init() {
	// Register available commands.
	Commands["reset-pwd"] = Command{
		Name:    "reset-pwd",
		Usage:   "./GoStreamRecord reset-pwd <username> <new-password>",
		Execute: resetPwdCommand,
	}
	Commands["add-user"] = Command{
		Name:    "reset-pwd",
		Usage:   "./GoStreamRecord add-user <username> <password>",
		Execute: addNewUser,
	}
}
