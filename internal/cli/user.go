package cli

import (
	"GoStreamRecord/internal/cli/color"
	"GoStreamRecord/internal/db"
	dbuser "GoStreamRecord/internal/db/users"
	"GoStreamRecord/internal/handlers/login"
)

func addNewUser(args []string) {
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

	db.CheckJson("users", "users.json", &db.Config.Users)
	// Loop over the users in the database to check if the user exists
	for _, u := range db.Config.Users.Users {
		if u.Name == username {

			color.Println("Bred", "User already exists!")
			return
		}
	}

	db.Config.Users.Users = append(db.Config.Users.Users, dbuser.Login{Name: username, Key: string(login.HashedPassword(newPassword))})
	// Save updated user configuration.
	db.Config.Update("users", "users.json", &db.Config.Users)

	color.Print("green", "Added new user ")
	color.Println("Bwhite", username)
}
