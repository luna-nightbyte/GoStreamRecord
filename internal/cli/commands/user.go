package commands

import (
	"GoStreamRecord/internal/cli/color"
	"GoStreamRecord/internal/db"
	dbuser "GoStreamRecord/internal/db/users"
	"GoStreamRecord/internal/web/handlers/login"
	"fmt"
)

func addNewUser(args []string) (string, error) {
	passwordWasReset = true

	if len(args) < 2 {
		// Provide clear feedback on what is missing.
		err := ""
		if len(args) < 1 {
			err = "No username provided."
			color.Println("Bred", err)
		} else {
			err = "No new password provided."
			color.Println("Bred", "No new password provided.")
		}
		PrintUsage(nil)
		return "", fmt.Errorf("%s", err)
	}

	username := args[0]
	newPassword := args[1]

	db.CheckJson("users", "users.json", &db.Config.Users)
	// Loop over the users in the database to check if the user exists
	for _, u := range db.Config.Users.Users {
		if u.Name == username {

			color.Println("Bred", "User already exists!")
			return "", fmt.Errorf("%s", "User already exists!")
		}
	}

	db.Config.Users.Users = append(db.Config.Users.Users, dbuser.Login{Name: username, Key: string(login.HashedPassword(newPassword))})
	// Save updated user configuration.
	db.Config.Update("users", "users.json", &db.Config.Users)

	color.Print("green", "Added new user ")
	color.Println("Bwhite", username)
	return "Added new user " + username, nil
}
