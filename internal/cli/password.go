package cli

import (
	"GoStreamRecord/internal/cli/color"
	"GoStreamRecord/internal/db"
	"GoStreamRecord/internal/handlers/login"
)

func resetPwdCommand(args []string) {
	passwordWasReset = true

	if len(args) < 2 {
		// Provide clear feedback on what is missing.
		if len(args) < 1 {
			color.Println("Bred", "No username provided.")
		} else {
			color.Println("Bred", "No new password provided.")
		}
		PrintUsage()
		return
	}

	username := args[0]
	newPassword := args[1]

	userFound := false

	// Loop over the users in the database to find a matching username.
	for i, u := range db.Config.Users.Users {
		if u.Name == username {
			db.Config.Users.Users[i].Key = string(login.HashedPassword(newPassword))
			userFound = true
			break
		}
	}

	if !userFound {
		color.Println("Bred", "No matching username found.")
		return
	}

	// Save updated user configuration.
	db.Config.Update("users", "users.json", &db.Config.Users)

	color.Print("green", "Password updated for ")
	color.Println("Bwhite", username)
}
