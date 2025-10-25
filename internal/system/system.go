package system

import (
	"context"
	"fmt"
	"log"
	"os"
	"remoteCtrl/internal/command"
	"remoteCtrl/internal/db"
	"remoteCtrl/internal/system/cookies"
	"remoteCtrl/internal/system/prettyprint"
	"sort"
)

func init() {
	// Register available commands.
	command.CMD.Startup.Add("reset-pwd", "./GoStreamRecord reset-pwd <username> <new-password>", ResetUserPassword)
	command.CMD.Startup.Add("add-user", "./GoStreamRecord add-user <username> <password> <role> <group-name>", AddNewUser)
	command.CMD.Startup.Add("del-user", "./GoStreamRecord del-user <username>", DeleteUser)
	command.CMD.Startup.Add("add-api", "./GoStreamRecord add-api <username> <api-name>", AddNewApi)
	command.CMD.Startup.Add("add-group", "./GoStreamRecord add-group <group-name> <description>", addGroup)
	command.CMD.Startup.Add("add-user-to-group", "./GoStreamRecord add-user-to-group <username> <group-name>", addUserToGroup)
	command.CMD.Startup.Add("list-users", "./GoStreamRecord list-users", listUsers)
	command.CMD.Startup.Add("list-groups", "./GoStreamRecord list-groups", listGroups)
	command.CMD.Startup.Add("list-roles", "./GoStreamRecord list-roles", listRoles)
	command.CMD.Startup.Add("list-user-groups", "./GoStreamRecord list-user-groups <username>", listUserGroups)
	command.CMD.Startup.Add("help", "./GoStreamRecord help", printUsage)
}

type Core struct {
	IsOnline        bool
	WaitForNetwork  bool
	Config          db.Config
	Context         context.Context
	Cancel          context.CancelFunc
	triggerShutdown chan os.Signal
}

var (
	System        Core
	onlineCheckIP string = "8.8.8.8"
	enableDebug   string = "false"
	DEBUG                = enableDebug == "true"
)

func StartupError() {
	PrintUsage()
}
func printUsage(args []string) {
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
func ResetUserPassword(args []string) {
	if len(args) < 2 {
		// Provide clear feedback on what is missing.
		if len(args) < 1 {
			prettyprint.P.LightRed.Println("No username provided.")
		} else {
			prettyprint.P.LightRed.Println("No new password provided.", args)
		}
		prettyprint.P.LightRed.Println("Error. See usage.")
		return
	}

	username := args[0]
	newPassword := args[1]

	userFound := false

	usrs, _ := db.DataBase.ListUsers()
	// Loop over the users in the database to find a matching username.
	for _, u := range usrs {
		if u.Username == username {
			db.DataBase.UpdateUser(u.ID, u.Username, string(cookies.HashedPassword(newPassword)))
			userFound = true
			break
		}
	}

	if !userFound {
		log.Println("No matching username found.")
		prettyprint.P.LightRed.Println("No matching username found.")
		return
	}

	log.Println("Password updated for", username)
	prettyprint.P.Success.Println(fmt.Sprintf("Password updated for %s", username))
}

// usrName, usrPass, UsrGroup, GroupRole
func AddNewUser(args []string) {
	switch len(args) {
	case 0:
		prettyprint.P.LightRed.Println("No username provided.")
		prettyprint.P.LightRed.Println("No new password provided.")
		prettyprint.P.LightRed.Println("No group name provided.")
		prettyprint.P.LightRed.Println("No group role provided.")
		return
	case 1:
		prettyprint.P.LightRed.Println("No new password provided.")
		prettyprint.P.LightRed.Println("No group name provided.")
		prettyprint.P.LightRed.Println("No group role provided.")
		return
	case 2:
		prettyprint.P.LightRed.Println("No group name provided.")
		prettyprint.P.LightRed.Println("No group role provided.")
		return
	case 3:
		prettyprint.P.LightRed.Println("No group role provided.")
		return
	case 4:
		break
	default:
		prettyprint.P.LightRed.Println("Too many arguments. See 'help' for usage")
		return
	}

	username := args[0]
	newPassword := args[1]
	role := args[2]
	group := args[3]

	err := db.DataBase.NewUser(username, newPassword)
	if err != nil {
		log.Println(err)
		prettyprint.P.LightRed.Println(err)
		return
	}

	user_id := db.DataBase.UserNameToID(username)
	group_id := db.DataBase.GroupNameToID(group)
	err = db.DataBase.AddUserToGroup(user_id, group_id, role)
	if err != nil {
		prettyprint.P.LightRed.Println(err)
		return
	}

	prettyprint.P.LightGreen.Print("Added new user: ")
	prettyprint.P.LightWhite.Println(username)
	prettyprint.P.LightGreen.Print("Password: ")
	prettyprint.P.FaintWhite.Println(newPassword)
	prettyprint.P.LightGreen.Print("Group: ")
	prettyprint.P.FaintWhite.Println(group)
	prettyprint.P.LightGreen.Print("Role: ")
	prettyprint.P.FaintWhite.Println(role)
}

func DeleteUser(args []string) {

	if len(args) < 1 {
		prettyprint.P.LightRed.Println("No username provided.")
		return
	}

	username := args[0]
	user_id := db.DataBase.UserNameToID(username)

	// Remove user from all groups
	groups, _, err := db.DataBase.ListGroupsByUserID(user_id)
	for _, group := range groups {
		err = db.DataBase.RemoveUserFromGroup(user_id, group.ID)
		if err != nil {
			log.Println(err)
			prettyprint.P.LightRed.Println(err)
			return
		}
	}

	// Delete all APIs for the user
	apis, err := db.DataBase.ListAvailableAPIsForUser(user_id)
	if err != nil {
		log.Println(err)
		prettyprint.P.LightRed.Println(err)
		return
	}
	for _, api := range apis {
		err = db.DataBase.DeleteApiForUser(user_id, api.ID)
		if err != nil {
			log.Println(err)
			prettyprint.P.LightRed.Println(err)
			return
		}
	}

	// Finally, delete the user
	err = db.DataBase.DeleteUser(user_id)
	if err != nil {
		log.Println(err)
		prettyprint.P.LightRed.Println(err)
		return
	}

	prettyprint.P.LightGreen.Print("Deleted user: ")
	prettyprint.P.LightWhite.Println(username)
}

func AddNewApi(args []string) {

	if len(args) < 2 {
		// Provide clear feedback on what is missing.
		if len(args) < 2 {
			prettyprint.P.LightRed.Println("No api provided.")
		} else if len(args) < 1 {
			prettyprint.P.LightRed.Println("No username provided.")
		}
		prettyprint.P.LightRed.Println("Error! See usage.")
		return
	}

	username := args[0]
	apiName := args[1]

	err := db.DataBase.NewApi(apiName, username)
	if err != nil {
		log.Println(err)
		prettyprint.P.LightRed.Println(err)
		return
	}
	user_id := db.DataBase.UserNameToID(username)
	apis, err := db.DataBase.ListAvailableAPIsForUser(user_id)
	if err != nil {
		prettyprint.P.LightRed.Println(err)
		return
	}

	prettyprint.P.LightGreen.Println("Added new api:")

	prettyprint.P.LightWhite.Println("KEY:", apis[apiName].Key)
	prettyprint.P.LightWhite.Println("Expires at:", apis[apiName].Expires)

}

// addGroup creates a new user group.
func addGroup(args []string) {
	if len(args) != 2 {
		log.Fatalf("Usage: ./GoStreamRecord add-group <group-name> <description>")
	}
	groupName := args[0]
	description := args[1]

	if err := db.DataBase.NewGroup(groupName, description); err != nil {
		log.Fatalf("Fatal: Could not create group '%s': %v", groupName, err)
	}

	prettyprint.P.LightGreen.Print("Successfully created group: ")
	prettyprint.P.LightWhite.Println(groupName)
	os.Exit(0)
}

// addUserToGroup assigns a user to an existing group.
func addUserToGroup(args []string) {
	if len(args) != 2 {
		log.Fatalf("Usage: ./GoStreamRecord add-user-to-group <username> <group-name>")
	}
	username := args[0]
	groupName := args[1]

	// Get User ID
	userID := db.DataBase.UserNameToID(username)
	if userID == 0 { // Assuming 0 is the "not found" indicator
		log.Fatalf("Fatal: User '%s' not found.", username)
	}

	// Get Group ID
	groupID := db.DataBase.GroupNameToID(groupName)
	if groupID == 0 { // Assuming 0 is the "not found" indicator
		log.Fatalf("Fatal: Group '%s' not found.", groupName)
	}

	// Add user to group with a default role
	// Assuming db.RoleUsers is the correct constant for a standard member
	if err := db.DataBase.AddUserToGroup(userID, groupID, db.RoleUsers); err != nil {
		log.Fatalf("Fatal: Could not add user '%s' to group '%s': %v", username, groupName, err)
	}

	fmt.Printf("Successfully added user %s to group %s\n", prettyprint.P.LightGreen.Color(username), prettyprint.P.LightGreen.Color(groupName))
	os.Exit(0)
}

// listUsers prints all registered users to the console.
func listUsers(args []string) {
	if len(args) != 0 {
		log.Fatalf("Usage: ./GoStreamRecord list-users")
	}

	users, err := db.DataBase.ListUsers() // Assumes returns map[string]db.User
	if err != nil {
		log.Fatalf("Fatal: Could not list users: %v", err)
	}

	prettyprint.P.LightWhite.Println("Registered Users:")
	if len(users) == 0 {
		prettyprint.P.FaintWhite.Println("  (No users found)")
		os.Exit(0)
	}

	// Get and sort usernames for clean output
	var userNames []string
	for name := range users {
		if name == db.InternalUser {
			continue
		}
		userNames = append(userNames, name)
	}
	sort.Strings(userNames)

	for _, name := range userNames {
		fmt.Printf("  - %s\n", prettyprint.P.LightGreen.Color(name))
	}
	os.Exit(0)
}

// listGroups prints all available groups to the console.
func listGroups(args []string) {
	if len(args) != 0 {
		log.Fatalf("Usage: ./GoStreamRecord list-groups")
	}

	groups, err := db.DataBase.ListGroups() // Assumes returns map[string]db.Group
	if err != nil {
		log.Fatalf("Fatal: Could not list groups: %v", err)
	}

	prettyprint.P.LightWhite.Println("Available Groups:")
	if len(groups) == 0 {
		prettyprint.P.FaintWhite.Println("  (No groups found)")
		os.Exit(0)
	}

	// Get and sort group names
	var groupNames []string
	for name := range groups {
		groupNames = append(groupNames, name)
	}
	sort.Strings(groupNames)

	for _, name := range groupNames {
		group := groups[name]
		fmt.Printf("  - %s (%s)\n", prettyprint.P.LightGreen.Color(group.Name), prettyprint.P.FaintWhite.Color(group.Description))
	}
	os.Exit(0)
}

// listUserGroups prints all groups a specific user belongs to.
func listUserGroups(args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: ./GoStreamRecord list-user-groups <username>")
	}
	username := args[0]

	userID := db.DataBase.UserNameToID(username)
	if userID == 0 {
		log.Fatalf("Fatal: User '%s' not found.", username)
	}
 
	groups, _, err := db.DataBase.ListGroupsByUserID(userID)
	if err != nil {
		log.Fatalf("Fatal: Could not list groups for user '%s': %v", username, err)
	}

	fmt.Printf("%s for user %s:\n", prettyprint.P.LightWhite.Color("Groups"), prettyprint.P.LightGreen.Color(username))
	if len(groups) == 0 {
		prettyprint.P.FaintWhite.Println("  (No groups assigned)")
		os.Exit(0)
	}

	for _, group := range groups {
		fmt.Printf("  - %s (%s)\n", prettyprint.P.LightGreen.Color(group.Name), prettyprint.P.FaintWhite.Color(group.Description))
	}
	os.Exit(0)
}

func listRoles(args []string) {
	if len(args) != 0 {
		log.Fatalf("Usage: ./GoStreamRecord list-roles")
	}

	usrs, _ := db.DataBase.ListUsers()
	for _, usr := range usrs {
		fmt.Println("Checking user id", usr)
		groupRelations, err := db.DataBase.GetUserGroupRelations(usr.ID) // Assumes returns map[string]db.Group
		if err != nil {
			fmt.Printf("Fatal: Could not list groups: %v", err)
			log.Fatalf("Fatal: Could not list groups: %v", err)
		}
		groups, _, _ := db.DataBase.ListGroupsByUserID(usr.ID)
		if err != nil {
			log.Fatalf("Fatal: Could not list groups: %v", err)
		}
		if len(groupRelations) == 0 {
			prettyprint.P.FaintWhite.Println(fmt.Sprintf("  (No group relations found for %s)", usr.Username))
			fmt.Println(groups)
			continue
		}
		for _, group := range groupRelations {
			if group.UserID == usr.ID {

				fmt.Printf(" | User: %s| - | Group: %s| Role: %s| Description: %s\n",
					prettyprint.P.Green.Color(usr.Username),
					prettyprint.P.LightGreen.Color(groups[group.GroupID].Name),
					prettyprint.P.FaintWhite.Color(group.Role),
					prettyprint.P.FaintWhite.Color(groups[group.GroupID].Description))
			}

		}
	}

	os.Exit(0)
}
