package cli

import (
	"GoStreamRecord/internal/cli/color"
	"GoStreamRecord/internal/db"
	"GoStreamRecord/internal/handlers/login"
	"fmt"
)

func resetPwdCommand(args []string) (string, error) {
	passwordWasReset = true

	if len(args) < 2 {
		// Provide clear feedback on what is missing.
		err := ""
		if len(args) < 1 {
			err = "No username provided."
			color.Println("Bred", err)
		} else {
			err = "No new password provided."
			color.Println("Bred", err)
		}
		PrintUsage(nil)
		return "", fmt.Errorf("%s", err)
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
		return "", fmt.Errorf("%s", "No matching username found.")
	}

	// Save updated user configuration.
	db.Config.Update("users", "users.json", &db.Config.Users)

	color.Print("green", "Password updated for ")
	color.Println("Bwhite", username)
	return "Password updated for " + username, nil
}
