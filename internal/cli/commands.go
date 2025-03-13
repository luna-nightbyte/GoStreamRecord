package cli

var passwordWasReset bool

// Command represents a CLI command with its name, usage, and execution function.
type Command struct {
	Name    string
	Usage   usageStruct
	Execute func(args []string)
}
type usageStruct struct {
	Bin     string
	Command string
	Args    string
}

// Global command registry.
var Commands = map[string]Command{}

func init() {
	// Register available commands.
	Commands["reset-pwd"] = Command{
		Name:    "reset-pwd",
		Usage:   usageStruct{Bin: "./GoStreamRecord", Command: "reset-pwd", Args: "<username> <new-password>"},
		Execute: resetPwdCommand,
	}
	Commands["add-user"] = Command{
		Name:    "reset-pwd",
		Usage:   usageStruct{Bin: "./GoStreamRecord", Command: "add-pwd", Args: "<username> <new-password>"},
		Execute: addNewUser,
	}
}
