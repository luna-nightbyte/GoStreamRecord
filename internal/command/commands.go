package command

import "fmt"

// Command represents a CLI command with its name, usage, and execution function.
type Command struct {
	Map map[string]commandStruct
}
type commandMap struct {
	Map map[string]commandStruct
}

// Command represents a CLI command with its name, usage, and execution function.
type commandStruct struct {
	Name  string
	Usage string
	Run   func(args []string)
}
type CMDs struct {
	Startup   Command
	Telemetry Command
	General   Command
}

var CMD CMDs

func init() {
	CMD.Startup.Map = make(map[string]commandStruct)
	CMD.Telemetry.Map = make(map[string]commandStruct)
}
func (c *Command) Add(commandName, usage string, runFunc func(args []string)) {
	if c.Map[commandName].Name != "" {
		fmt.Println("Command already in use.", c.Map[commandName])
		return
	}
	c.Map[commandName] = commandStruct{
		Name:  commandName,
		Usage: usage,
		Run:   runFunc,
	}
}
