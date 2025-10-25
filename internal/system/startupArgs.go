package system

import (
	"remoteCtrl/internal/command"
	"remoteCtrl/internal/db"
	"remoteCtrl/internal/system/prettyprint"
	"sort"
)

func init() {
	// Register available commands.
	command.CMD.Startup.Add("reset-pwd", "./GoStreamRecord reset-pwd <username> <new-password>", db.ResetUserPassword)
	command.CMD.Startup.Add("add-user", "./GoStreamRecord add-user <username> <password> <role> <group-name>", db.AddNewUser)
	command.CMD.Startup.Add("del-user", "./GoStreamRecord del-user <username>", db.DeleteUser)
	command.CMD.Startup.Add("add-api", "./GoStreamRecord add-api <username> <api-name>", db.AddNewApi)
	command.CMD.Startup.Add("add-group", "./GoStreamRecord add-group <group-name> <description>", db.AddGroup)
	command.CMD.Startup.Add("add-user-to-group", "./GoStreamRecord add-user-to-group <username> <group-name>", db.AddUserToGroup)
	command.CMD.Startup.Add("list-users", "./GoStreamRecord list-users", db.ListUsers)
	command.CMD.Startup.Add("list-groups", "./GoStreamRecord list-groups", db.ListGroups)
	command.CMD.Startup.Add("list-roles", "./GoStreamRecord list-roles", db.ListRoles)
	command.CMD.Startup.Add("list-user-groups", "./GoStreamRecord list-user-groups <username>", db.ListUserGroups)
	command.CMD.Startup.Add("help", "./GoStreamRecord help", printUsage)
}
func StartupError() {
	PrintUsage()
}
func printUsage(args ...string) {
	PrintUsage()

}
func PrintUsage() {
	keys := make([]string, 0, len(command.CMD.Startup.Map))
	for key := range command.CMD.Startup.Map {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	prettyprint.P.LightCyan.Println("Usage:")

	for _, key := range keys {
		cmd := command.CMD.Startup.Map[key]
		prettyprint.P.LightCyan.Println(" - " + cmd.Usage)
	}

	prettyprint.P.LightCyan.Println("Otherwise run the server without any arguments.")
}
