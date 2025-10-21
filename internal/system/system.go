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
	"remoteCtrl/internal/system/settings"
	"sort"
)

func init() {
	// Register available commands.
	command.CMD.Startup.Add("reset-pwd", "./GoStreamRecord reset-pwd <username> <new-password>", ResetWebUIPassword)
	command.CMD.Startup.Add("add-user", "./GoStreamRecord add-user <username> <password>", AddNewUser)
	// --- User and Group Management Commands ---
	command.CMD.Startup.Add("add-group", "./GoStreamRecord add-group <group-name> <description>", addGroup)
	command.CMD.Startup.Add("add-user-to-group", "./GoStreamRecord add-user-to-group <username> <group-name>", addUserToGroup)
	command.CMD.Startup.Add("list-users", "./GoStreamRecord list-users", listUsers)
	command.CMD.Startup.Add("list-groups", "./GoStreamRecord list-groups", listGroups)
	command.CMD.Startup.Add("list-user-groups", "./GoStreamRecord list-user-groups <username>", listUserGroups)
	command.CMD.Startup.Add("help", "./GoStreamRecord help", printUsage)
}

type Core struct {
	IsOnline        bool
	WaitForNetwork  bool
	Config          settings.DB // ./settings/settings.json
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

func Init() error {
	return nil
}

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

	fmt.Println(prettyprint.Cyan("Usage:"))

	for _, key := range keys {
		cmd := command.CMD.Startup.Map[key]
		fmt.Println(prettyprint.Cyan(" - " + cmd.Usage))
	}

	fmt.Println(prettyprint.Cyan("Otherwise run the server without any arguments."))
}
func ResetWebUIPassword(args []string) {
	if len(args) < 2 {
		// Provide clear feedback on what is missing.
		if len(args) < 1 {
			fmt.Println(prettyprint.BoldRed("No username provided."))
		} else {
			fmt.Println(prettyprint.BoldRed("No new password provided.", args))
		}
		fmt.Println(prettyprint.BoldRed("Error. See usage."))
		return
	}

	username := args[0]
	newPassword := args[1]

	userFound := false

	usrs, _ := db.DataBase.Users.List()
	// Loop over the users in the database to find a matching username.
	for _, u := range usrs {
		if u.Username == username {
			db.DataBase.Users.Update(u.ID, u.Username, string(cookies.HashedPassword(newPassword)))
			userFound = true
			break
		}
	}

	if !userFound {
		log.Println("No matching username found.")
		fmt.Println(prettyprint.BoldRed("No matching username found."))
		return
	}

	log.Println("Password updated for", username)
	fmt.Println(prettyprint.Green("Password updated for "), prettyprint.BoldWhite(username))
}

func AddNewUser(args []string) {

	if len(args) < 3 {
		// Provide clear feedback on what is missing.
		if len(args) < 2 {
			fmt.Println(prettyprint.BoldRed("No role provided."))
		} else if len(args) < 1 {
			fmt.Println(prettyprint.BoldRed("No username provided."))
		} else {
			fmt.Println(prettyprint.BoldRed("No new password provided.", args))
		}
		fmt.Println(prettyprint.BoldRed("Error! See usage."))
		return
	}

	username := args[0]
	newPassword := args[1]

	role := args[2]
	group := args[3]

	err := db.DataBase.Users.New(username, newPassword)
	if err != nil {
		log.Println(err)
		fmt.Println(prettyprint.BoldRed(err))
		return
	}

	user_id := db.DataBase.Users.NameToID(username)
	group_id := db.DataBase.Groups.NameToID(group)
	err = db.DataBase.Groups.AddUser(user_id, group_id, role)
	if err != nil {
		fmt.Println(prettyprint.BoldRed(err))
		return
	}

	fmt.Println(prettyprint.Green("Added new user:"), prettyprint.BoldWhite(username))
	fmt.Println(prettyprint.Green("Password:"), prettyprint.BoldGrey(newPassword))
	fmt.Println(prettyprint.Green("Group:"), prettyprint.BoldGrey(group))
	fmt.Println(prettyprint.Green("Role:"), prettyprint.BoldGrey(role))
}

// addGroup creates a new user group.
func addGroup(args []string) {
	if len(args) != 2 {
		log.Fatalf("Usage: ./GoStreamRecord add-group <group-name> <description>")
	}
	groupName := args[0]
	description := args[1]

	if err := db.DataBase.Groups.New(groupName, description); err != nil {
		log.Fatalf("Fatal: Could not create group '%s': %v", groupName, err)
	}

	fmt.Printf("Successfully created group: %s\n", prettyprint.Green(groupName))
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
	userID := db.DataBase.Users.NameToID(username)
	if userID == 0 { // Assuming 0 is the "not found" indicator
		log.Fatalf("Fatal: User '%s' not found.", username)
	}

	// Get Group ID
	groupID := db.DataBase.Groups.NameToID(groupName)
	if groupID == 0 { // Assuming 0 is the "not found" indicator
		log.Fatalf("Fatal: Group '%s' not found.", groupName)
	}

	// Add user to group with a default role
	// Assuming db.RoleUsers is the correct constant for a standard member
	if err := db.DataBase.Groups.AddUser(userID, groupID, db.RoleUsers); err != nil {
		log.Fatalf("Fatal: Could not add user '%s' to group '%s': %v", username, groupName, err)
	}

	fmt.Printf("Successfully added user %s to group %s\n", prettyprint.Green(username), prettyprint.Green(groupName))
	os.Exit(0)
}

// listUsers prints all registered users to the console.
func listUsers(args []string) {
	if len(args) != 0 {
		log.Fatalf("Usage: ./GoStreamRecord list-users")
	}

	users, err := db.DataBase.Users.List() // Assumes returns map[string]db.User
	if err != nil {
		log.Fatalf("Fatal: Could not list users: %v", err)
	}

	fmt.Println(prettyprint.BoldWhite("Registered Users:"))
	if len(users) == 0 {
		fmt.Println(prettyprint.BoldGrey("  (No users found)"))
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
		fmt.Printf("  - %s\n", prettyprint.Green(name))
	}
	os.Exit(0)
}

// listGroups prints all available groups to the console.
func listGroups(args []string) {
	if len(args) != 0 {
		log.Fatalf("Usage: ./GoStreamRecord list-groups")
	}

	groups, err := db.DataBase.Groups.List() // Assumes returns map[string]db.Group
	if err != nil {
		log.Fatalf("Fatal: Could not list groups: %v", err)
	}

	fmt.Println(prettyprint.BoldWhite("Available Groups:"))
	if len(groups) == 0 {
		fmt.Println(prettyprint.BoldGrey("  (No groups found)"))
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
		// Assumes db.Group struct has Name and Description fields
		fmt.Printf("  - %s (%s)\n", prettyprint.Green(group.Name), prettyprint.BoldGrey(group.Description))
	}
	os.Exit(0)
}

// listUserGroups prints all groups a specific user belongs to.
func listUserGroups(args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: ./GoStreamRecord list-user-groups <username>")
	}
	username := args[0]

	userID := db.DataBase.Users.NameToID(username)
	if userID == 0 {
		log.Fatalf("Fatal: User '%s' not found.", username)
	}

	// This function is ASSUMED to exist: db.DataBase.Groups.ListForUser(userID)
	// It's assumed to return a slice of db.Group structs.
	groups, _, err := db.DataBase.Groups.ListGroupsByUserID(userID)
	if err != nil {
		log.Fatalf("Fatal: Could not list groups for user '%s': %v", username, err)
	}

	fmt.Printf("%s for user %s:\n", prettyprint.BoldWhite("Groups"), prettyprint.Green(username))
	if len(groups) == 0 {
		fmt.Println(prettyprint.BoldGrey("  (No groups assigned)"))
		os.Exit(0)
	}

	for _, group := range groups {
		fmt.Printf("  - %s (%s)\n", prettyprint.Green(group.Name), prettyprint.BoldGrey(group.Description))
	}
	os.Exit(0)
}
