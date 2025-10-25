package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"remoteCtrl/internal/system/prettyprint"
	"remoteCtrl/internal/utils"
	"sort"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func (db *DB) queryApiSql(query string, args ...any) (Api, error) {
	var a Api
	row := db.SQL.QueryRowContext(db.ctx, query, args...)
	err := row.Scan(&a.ID, &a.Name, &a.Key, &a.Expires, &a.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return a, ErrNotFound
		}
		return a, err
	}

	return a, nil
}
func (db *DB) queryGroupSql(query string, args ...any) (Group, error) {
	row := db.SQL.QueryRowContext(db.ctx, query, args...)
	var group Group
	var createdAt string
	err := row.Scan(&group.ID, &group.Name, &group.Description, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return group, ErrNotFound
		}
		return group, err
	}

	if group.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
		return group, err
	}
	return group, nil
}
func (db *DB) queryStreamerSql(query string, args ...any) (Streamer, error) {
	var s Streamer
	row := db.SQL.QueryRowContext(db.ctx, query, args...)
	err := row.Scan(&s.ID, &s.Name, &s.Provider)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s, ErrNotFound
		}
		return s, err
	}

	return s, nil
}

func (db *DB) queryStreamerGroupRelationsSql(query string, args ...any) (streamer_group_relations, error) {
	row := db.SQL.QueryRowContext(db.ctx, query, args...)
	var streamerGrp streamer_group_relations
	err := row.Scan(&streamerGrp.StreamerID, &streamerGrp.GroupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return streamerGrp, ErrUserNotFound
		}
		return streamerGrp, err
	}

	return streamerGrp, nil
}
func (db *DB) queryTabSql(query string, args ...any) (Tab, error) {
	var t Tab
	row := db.SQL.QueryRowContext(db.ctx, query, args...)
	err := row.Scan(&t.ID, &t.Name, &t.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return t, ErrNotFound
		}
		return t, err
	}

	return t, nil
}

func (db *DB) queryTabGroupRelationsSql(query string, args ...any) (tab_group_relations, error) {
	row := db.SQL.QueryRowContext(db.ctx, query, args...)
	var usrGrp tab_group_relations
	err := row.Scan(&usrGrp.TabID, &usrGrp.GroupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return usrGrp, ErrUserNotFound
		}
		return usrGrp, err
	}

	return usrGrp, nil
}

// username, newPassword
func ResetUserPassword(args ...string) {
	switch len(args) {
	case 0:
		prettyprint.P.LightRed.Println("No username provided.")
		prettyprint.P.LightRed.Println("No new password provided.")
		return
	case 1:
		prettyprint.P.LightRed.Println("No new password provided.")
		return
	case 2:
		break
	default:
		prettyprint.P.LightRed.Println("Too many arguments. See 'help' for usage")
		return
	}
	username := args[0]
	newPassword := args[1]

	userFound := false

	usrs, _ := DataBase.ListUsers()
	// Loop over the users in the database to find a matching username.
	for _, u := range usrs {
		if u.Username == username {
			hash, _ := utils.HashPassword(newPassword)
			DataBase.UpdateUser(u.ID, u.Username, string(hash))
			userFound = true
			break
		}
	}

	if !userFound {
		log.Println("No matching username found.")
		prettyprint.P.LightRed.Println("No matching username found.")
		return
	}

	log.Println("Password updated for", username)
	prettyprint.P.Success.Println(fmt.Sprintf("Password updated for %s", username))
}

// usrName, usrPass, UsrGroup, GroupRole
func AddNewUser(args ...string) {
	switch len(args) {
	case 0:
		prettyprint.P.LightRed.Println("No username provided.")
		prettyprint.P.LightRed.Println("No new password provided.")
		prettyprint.P.LightRed.Println("No group name provided.")
		prettyprint.P.LightRed.Println("No group role provided.")
		return
	case 1:
		prettyprint.P.LightRed.Println("No new password provided.")
		prettyprint.P.LightRed.Println("No group name provided.")
		prettyprint.P.LightRed.Println("No group role provided.")
		return
	case 2:
		prettyprint.P.LightRed.Println("No group name provided.")
		prettyprint.P.LightRed.Println("No group role provided.")
		return
	case 3:
		prettyprint.P.LightRed.Println("No group role provided.")
		return
	case 4:
		break
	default:
		prettyprint.P.LightRed.Println("Too many arguments. See 'help' for usage")
		return
	}

	username := args[0]
	newPassword := args[1]
	role := args[2]
	group := args[3]

	err := DataBase.NewUser(username, newPassword)
	if err != nil {
		log.Println(err)
		prettyprint.P.LightRed.Println(err)
		return
	}

	user_id := DataBase.UserNameToID(username)

	group_id := DataBase.GroupNameToID(group)
	if group_id == -1 {
		prettyprint.P.BoldRed.Println("Group does not exist")
		return
	}

	err = DataBase.AddUserToGroup(user_id, group_id, role)
	if err != nil {
		prettyprint.P.LightRed.Println(err)
		return
	}

	prettyprint.P.LightGreen.Print("Added new user: ")
	prettyprint.P.LightWhite.Println(username)
	prettyprint.P.LightGreen.Print("Password: ")
	prettyprint.P.FaintWhite.Println(newPassword)
	prettyprint.P.LightGreen.Print("Group: ")
	prettyprint.P.FaintWhite.Println(group)
	prettyprint.P.LightGreen.Print("Role: ")
	prettyprint.P.FaintWhite.Println(role)
}

func DeleteUser(args ...string) {
	switch len(args) {
	case 0:
		prettyprint.P.LightRed.Println("No username provided.")
		return
	case 1:
		break
	default:
		prettyprint.P.LightRed.Println("Too many arguments. See 'help' for usage")
		return
	}
	username := args[0]
	user_id := DataBase.UserNameToID(username)

	// Remove user from all groups
	groups, _, err := DataBase.ListGroupsByUserID(user_id)
	for _, group := range groups {
		err = DataBase.RemoveUserFromGroup(user_id, group.ID)
		if err != nil {
			log.Println(err)
			prettyprint.P.LightRed.Println(err)
			return
		}
	}

	// Delete all APIs for the user
	apis, err := DataBase.ListAvailableAPIsForUser(user_id)
	if err != nil {
		log.Println(err)
		prettyprint.P.LightRed.Println(err)
		return
	}
	for _, api := range apis {
		err = DataBase.DeleteApiForUser(user_id, api.ID)
		if err != nil {
			log.Println(err)
			prettyprint.P.LightRed.Println(err)
			return
		}
	}

	// Finally, delete the user
	err = DataBase.DeleteUser(user_id)
	if err != nil {
		log.Println(err)
		prettyprint.P.LightRed.Println(err)
		return
	}

	prettyprint.P.LightGreen.Print("Deleted user: ")
	prettyprint.P.LightWhite.Println(username)
}

func AddNewApi(args ...string) {
	switch len(args) {
	case 0:
		prettyprint.P.LightRed.Println("No username provided.")
		prettyprint.P.LightRed.Println("No new api provided.")
		return
	case 1:
		prettyprint.P.LightRed.Println("No new api provided.")
		return
	case 2:
		break
	default:
		prettyprint.P.LightRed.Println("Too many arguments. See 'help' for usage")
		return
	}
	username := args[0]
	apiName := args[1]

	user, err := DataBase.GetUserByName(username)
	if err != nil {
		prettyprint.P.LightRed.Println(err)
		return
	}
	err = DataBase.NewApi(apiName, user)
	if err != nil {
		if err == ErrNotFound {
			prettyprint.P.LightRed.Println(ErrUserNotFound)
			return
		}
		fmt.Println("Error creating new api")
		prettyprint.P.LightRed.Println(err)
		return
	}
	apis, err := DataBase.ListAvailableAPIsForUser(user.ID)
	if err != nil {
		prettyprint.P.LightRed.Println(err)
		return
	}

	prettyprint.P.LightGreen.Println("Added new api:")

	prettyprint.P.LightWhite.Println("KEY:", apis[apiName].Key)
	prettyprint.P.LightWhite.Println("Expires at:", apis[apiName].Expires)

}

// addGroup creates a new user group.
func AddGroup(args ...string) {
	var description string

	switch len(args) {
	case 0:
		prettyprint.P.LightRed.Println("No group provided.")
		prettyprint.P.LightRed.Println("No description provided.")
		return
	case 1:
		break // empty description
	case 2:
		description = args[1]
	default:
		prettyprint.P.LightRed.Println("Too many arguments. See 'help' for usage")
		return
	}
	groupName := args[0]

	if err := DataBase.NewGroup(groupName, description); err != nil {
		log.Fatalf("Fatal: Could not create group '%s': %v", groupName, err)
	}

	prettyprint.P.LightGreen.Print("Successfully created group: ")
	prettyprint.P.LightWhite.Println(groupName)
	os.Exit(0)
}

// addUserToGroup assigns a user to an existing group.
func AddUserToGroup(args ...string) {

	switch len(args) {
	case 0:
		prettyprint.P.LightRed.Println("No name provided.")
		prettyprint.P.LightRed.Println("No group provided.")
		prettyprint.P.LightRed.Println("No role provided.")
		return
	case 1:
		prettyprint.P.LightRed.Println("No group provided.")
		prettyprint.P.LightRed.Println("No role provided.")
		return
		return
	case 2:
		prettyprint.P.LightRed.Println("No role provided.")
		return
	case 3:
		break
	default:
		prettyprint.P.LightRed.Println("Too many arguments. See 'help' for usage")
		return
	}
	username := args[0]
	groupName := args[1]
	role := args[2]

	// Get User ID
	userID := DataBase.UserNameToID(username)
	if userID == 0 { // Assuming 0 is the "not found" indicator
		log.Fatalf("Fatal: User '%s' not found.", username)
	}

	// Get Group ID
	groupID := DataBase.GroupNameToID(groupName)
	if groupID == 0 { // Assuming 0 is the "not found" indicator
		log.Fatalf("Fatal: Group '%s' not found.", groupName)
	}

	if role != RoleAdmin && role != RoleUsers {
		log.Fatalf("Fatal: Invalid role '%s'\nUse either '%s' or '%s'", groupName, RoleUsers, RoleAdmin)
	}
	// Add user to group with a default role
	// Assuming RoleUsers is the correct constant for a standard member
	if err := DataBase.AddUserToGroup(userID, groupID, role); err != nil {
		log.Fatalf("Fatal: Could not add user '%s' to group '%s': %v", username, groupName, err)
	}
	if username == InternalUser {
		return
	}

	fmt.Printf("Added user %s to group %s\n", prettyprint.P.LightGreen.Color(username), prettyprint.P.LightGreen.Color(groupName))

}

// listUsers prints all registered users to the console.
func ListUsers(args ...string) {
	users, err := DataBase.ListUsers() // Assumes returns map[string]User
	if err != nil {
		log.Fatalf("Fatal: Could not list users: %v", err)
	}

	prettyprint.P.LightWhite.Println("Registered Users:")
	if len(users) == 0 {
		prettyprint.P.FaintWhite.Println("  (No users found)")
		os.Exit(0)
	}

	// Get and sort usernames for clean output
	var userNames []string
	for name := range users {
		if name == InternalUser {
			continue
		}
		userNames = append(userNames, name)
	}
	sort.Strings(userNames)

	for _, name := range userNames {
		fmt.Printf("  - %s\n", prettyprint.P.LightGreen.Color(name))
	}
	os.Exit(0)
}

// listGroups prints all available groups to the console.
func ListGroups(args ...string) {

	groups, err := DataBase.ListAllGroups() // Assumes returns map[string]Group
	if err != nil {
		log.Fatalf("Fatal: Could not list groups: %v", err)
	}

	prettyprint.P.LightWhite.Println("Available Groups:")
	if len(groups) == 0 {
		prettyprint.P.FaintWhite.Println("  (No groups found)")
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
		fmt.Printf("  - %s (%s)\n", prettyprint.P.LightGreen.Color(group.Name), prettyprint.P.FaintWhite.Color(group.Description))
	}
	os.Exit(0)
}

// listUserGroups prints all groups a specific user belongs to.
func ListUserGroups(args ...string) {
	if len(args) != 1 {
		log.Fatalf("Usage: ./GoStreamRecord list-user-groups <username>")
	}
	username := args[0]

	userID := DataBase.UserNameToID(username)
	if userID == 0 {
		log.Fatalf("Fatal: User '%s' not found.", username)
	}

	groups, _, err := DataBase.ListGroupsByUserID(userID)
	if err != nil {
		log.Fatalf("Fatal: Could not list groups for user '%s': %v", username, err)
	}

	fmt.Printf("%s for user %s:\n", prettyprint.P.LightWhite.Color("Groups"), prettyprint.P.LightGreen.Color(username))
	if len(groups) == 0 {
		prettyprint.P.FaintWhite.Println("  (No groups assigned)")
		os.Exit(0)
	}

	for _, group := range groups {
		fmt.Printf("  - %s (%s)\n", prettyprint.P.LightGreen.Color(group.Name), prettyprint.P.FaintWhite.Color(group.Description))
	}
	os.Exit(0)
}

func ListRoles(args ...string) {
	groupRelations, err := DataBase.ListAllUserGroupRelations() // Assumes returns map[string]Group
	if err != nil {
		fmt.Printf("Fatal: Could not list groups: %v", err)
		log.Fatalf("Fatal: Could not list groups: %v", err)
	}
	if len(groupRelations) == 0 {
		prettyprint.P.FaintWhite.Println(fmt.Sprintf("No group relations found for"))
		return
	}

	if len(groupRelations) == 1 {
		if len(groupRelations[0].Role) == 0 {
			prettyprint.P.BoldRed.Println("No relations found..")
			return
		}
	}

	prettyprint.P.Green.Println("Users and roles:")
	for _, ugr := range groupRelations {
		if ugr.Role == "" {
			if len(groupRelations) == 1 {
				prettyprint.P.BoldRed.Println("No relations found..")
			}
			continue
		}
		user, _ := DataBase.GetUserByID(ugr.UserID)
		group, _ := DataBase.GetGroupByID(ugr.GroupID)

		fmt.Printf("  ____________________________\n | User: %s\n | Group: %s\n | Role: %s\n | Description: %s\n  ____________________________\n",
			prettyprint.P.Green.Color(user.Username),
			prettyprint.P.LightGreen.Color(group.Name),
			prettyprint.P.FaintWhite.Color(ugr.Role),
			prettyprint.P.FaintWhite.Color(group.Description))

	}

}
