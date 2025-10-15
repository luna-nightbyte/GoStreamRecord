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
	command.CMD.Startup.Add("reset-pwd", "./remoteCtrl reset-pwd <username> <new-password>", ResetWebUIPassword)
	command.CMD.Startup.Add("add-user", "./remoteCtrl add-user <username> <password>", AddNewUser)
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
			fmt.Println(prettyprint.BoldRed("No new password provided."))
		}
		fmt.Println(prettyprint.BoldRed("Error. See usage."))
		return
	}

	username := args[0]
	newPassword := args[1]

	userFound := false

	// Loop over the users in the database to find a matching username.
	for i, u := range System.DB.Users.Users {
		fmt.Println(username, u.Name)
		if u.Name == username {
			System.DB.Users.Users[i].Key = string(cookies.HashedPassword(newPassword))
			userFound = true
			break
		}
	}

	if !userFound {
		log.Println("No matching username found.")
		fmt.Println(prettyprint.BoldRed("No matching username found."))
		return
	}

	// Save updated user configuration.
	db.Update(settings.CONFIG_USERS_PATH, &System.DB.Users)
	log.Println("Password updated for", username)
	fmt.Println(prettyprint.Green("Password updated for "), prettyprint.BoldWhite(username))
}

func AddNewUser(args []string) {
	if len(args) < 2 {
		// Provide clear feedback on what is missing.
		if len(args) < 1 {
			fmt.Println(prettyprint.BoldRed("No username provided."))
		} else {
			fmt.Println(prettyprint.BoldRed("No new password provided."))
		}
		fmt.Println(prettyprint.BoldRed("Error! See usage."))
		return
	}

	username := args[0]
	newPassword := args[1]

	db.LoadConfig(settings.CONFIG_USERS_PATH, &System.DB.Users)
	// Loop over the users in the database to check if the user exists
	for _, u := range System.DB.Users.Users {
		fmt.Println(username, u.Name)
		if u.Name == username {
			log.Println("User already exists!")
			fmt.Println(prettyprint.BoldRed("User already exists!"))
			return
		}
	}

	System.DB.Users.Users = append(System.DB.Users.Users, settings.Login{Name: username, Key: string(cookies.HashedPassword(newPassword))})
	// Save updated user configuration.
	db.Update(settings.CONFIG_USERS_PATH, &System.DB.Users)
	log.Println("Added new user", username)
	fmt.Println(prettyprint.Green("Added new user "), prettyprint.BoldWhite(username))
}
