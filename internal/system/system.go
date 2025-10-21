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
)

type Core struct {
	IsOnline        bool
	WaitForNetwork  bool
	DB              settings.DB // ./settings/settings.json
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

func init() {
	// Register available commands.
	command.CMD.Startup.Add("reset-pwd", "./GoStreamRecord reset-pwd <username> <new-password>", ResetWebUIPassword)
	command.CMD.Startup.Add("add-user", "./GoStreamRecord add-user <username> <password>", AddNewUser)
}

func PrintUsage() {
	fmt.Println(prettyprint.Cyan("Usage:"))
	for _, cmd := range command.CMD.Startup.Map {
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

	fmt.Println("num of args:", len(args))

	fmt.Println("args:", (args))
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

	log.Println("Added new user", username)
	fmt.Println(prettyprint.Green("Added new user "), prettyprint.BoldWhite(username))
}
