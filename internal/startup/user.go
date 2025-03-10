package startup

import (
	"GoStreamRecord/internal/db"
	dbuser "GoStreamRecord/internal/db/users"
	"GoStreamRecord/internal/handlers/login"
	"GoStreamRecord/internal/prettyprint"
	"fmt"
	"log"
)

func addNewUser(args []string) {
	passwordWasReset = true

	if len(args) < 2 {
		// Provide clear feedback on what is missing.
		if len(args) < 1 {
			fmt.Println(prettyprint.BoldRed("No username provided."))
		} else {
			fmt.Println(prettyprint.BoldRed("No new password provided."))
		}
		fmt.Println(prettyprint.BoldRed("Usage:"), Commands["add-user"].Usage)
		return
	}

	username := args[0]
	newPassword := args[1]

	db.CheckJson("users", "users.json", &db.Config.Users)
	// Loop over the users in the database to check if the user exists
	for _, u := range db.Config.Users.Users {
		fmt.Println(username, u.Name)
		if u.Name == username {
			log.Println("User already exists!")
			fmt.Println(prettyprint.BoldRed("User already exists!"))
			return
		}
	}

	db.Config.Users.Users = append(db.Config.Users.Users, dbuser.Login{Name: username, Key: string(login.HashedPassword(newPassword))})
	// Save updated user configuration.
	db.Config.Update("users", "users.json", &db.Config.Users)
	log.Println("Added new user", username)
	fmt.Println(prettyprint.Green("Added new user "), prettyprint.BoldWhite(username))
}
