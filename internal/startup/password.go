package startup

import (
	"GoStreamRecord/internal/db"
	"GoStreamRecord/internal/handlers/login"
	"GoStreamRecord/internal/prettyprint"
	"fmt"
	"log"
)

func resetPwdCommand(args []string) {
	passwordWasReset = true

	if len(args) < 2 {
		// Provide clear feedback on what is missing.
		if len(args) < 1 {
			fmt.Println(prettyprint.BoldRed("No username provided."))
		} else {
			fmt.Println(prettyprint.BoldRed("No new password provided."))
		}
		fmt.Println(prettyprint.BoldRed("Usage:"), Commands["reset-pwd"].Usage)
		return
	}

	username := args[0]
	newPassword := args[1]

	userFound := false

	// Loop over the users in the database to find a matching username.
	for i, u := range db.Config.Users.Users {
		fmt.Println(username, u.Name)
		if u.Name == username {
			db.Config.Users.Users[i].Key = string(login.HashedPassword(newPassword))
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
	db.Config.Update("users", "users.json", &db.Config.Users)
	log.Println("Password updated for", username)
	fmt.Println(prettyprint.Green("Password updated for "), prettyprint.BoldWhite(username))
}
