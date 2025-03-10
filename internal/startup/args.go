package startup

import (
	"GoStreamRecord/internal/db"
	dbuser "GoStreamRecord/internal/db/users"
	"GoStreamRecord/internal/handlers/login"
	"GoStreamRecord/internal/prettyprint"
	"fmt"
	"log"
)

var passwordWasReset bool

// Command represents a CLI command with its name, usage, and execution function.
type Command struct {
	Name    string
	Usage   string
	Execute func(args []string)
}

// Global command registry.
var Commands = map[string]Command{}

func init() {
	// Register available commands.
	Commands["reset-pwd"] = Command{
		Name:    "reset-pwd",
		Usage:   "./GoStreamRecord reset-pwd <username> <new-password>",
		Execute: resetPwdCommand,
	}
	// Register available commands.
	Commands["add-user"] = Command{
		Name:    "reset-pwd",
		Usage:   "./GoStreamRecord add-user <username> <password>",
		Execute: addNewUser,
	}
}

func main() {

}

func PrintUsage() {
	fmt.Println(prettyprint.Cyan("Usage:"))
	for _, cmd := range Commands {
		fmt.Println(prettyprint.Cyan(" - " + cmd.Usage))
	}
	fmt.Println(prettyprint.Cyan("Otherwise run the server without any arguments."))
}

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
