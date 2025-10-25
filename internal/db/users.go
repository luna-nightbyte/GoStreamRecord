package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"remoteCtrl/internal/system/prettyprint"
	"remoteCtrl/internal/utils"
	"remoteCtrl/internal/web/handlers/cookie"
	"strings"
	"time"
)

// ErrUserNotFound is returned when a user is not found in the database.
var ErrUserNotFound = errors.New("user not found")
var ErrNotFound = errors.New("not found")

const ErrIsExist = "UNIQUE constraint failed"

// SQL QUERIES ---------------------------------------------------------------------

// AddUser hashes the password and inserts a new user record.
func (db *DB) NewUser(username, raw_password string) error {
	if username == "" || raw_password == "" {
		return errors.New("username and password cannot be empty")
	}
	hash, err := utils.HashPassword(raw_password)
	if err != nil {
		return fmt.Errorf("failed to hash pasd.ctx, csword: %w", err)
	}
	now := time.Now().Format(time.RFC3339)
	_, err = DataBase.SQL.ExecContext(db.ctx, createUser, username, hash, now)
	if err != nil {
		if strings.Contains(err.Error(), ErrIsExist) {
			return errors.New("Username already exists")
		}
		return err
	}

	return nil
}

// UpdateUser updates an existing user's details.
// If newPassword is an empty string, the password is not updated.
func (db *DB) UpdateUser(userID int, newUsername string, newPassword string) error {
	if newUsername == "" {
		return errors.New("username cannot be empty")
	}

	var result sql.Result
	var err error
	if newPassword != "" {
		hash, err := utils.HashPassword(newPassword)
		if err != nil {
			return fmt.Errorf("failed to hash new password: %w", err)
		}
		result, err = db.SQL.ExecContext(db.ctx, updateUser, newUsername, hash, userID)
	} else {
		usrs, _ := db.ListUsers()
		for _, urs := range usrs {
			if urs.ID == userID {
				result, err = db.SQL.ExecContext(db.ctx, updateUser, newUsername, urs.PasswordHash, userID)

			}
		}
		//query := "UPDATE users SET usernamename=?, group_ids=?, updated_at=? WHERE id=?"
	}

	if err != nil {
		return errors.New("failed to update user (username may already exist)")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// DeleteUser removes a user record by ID.
func (db *DB) DeleteUser(userID int) error {
	result, err := db.SQL.ExecContext(db.ctx, deleteUser, userID)
	if err != nil {
		return fmt.Errorf("database error during deletion: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// Authenticate checks a user's credentials against the database.
func (db *DB) AuthenticateUser(username, password string) (bool, error) {
	user, err := db.GetUserByName(username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return false, errors.New("invalid username or password")
		}
		return false, err
	}

	if utils.CheckPasswordHash(password, user.PasswordHash) {
		return true, nil
	}

	return false, errors.New("invalid username or password")
}

// IsAdmin checks ifunc (db *DB) Authenticate(username, password string) (bool, error)f a user has admin privileges.
func (db *DB) IsUserAdmin(username string) (bool, error) {
	user, err := db.GetUserByName(username)
	if err != nil {
		return false, err
	}
	_, role, err := db.ListGroupsByUserID(user.ID)
	if role == RoleAdmin {
		return true, nil
	}
	return false, err
}

func (db *DB) RequestUserID(r *http.Request) int {
	name, err := cookie.ValidateSession(r)
	if err != nil {
		prettyprint.P.Error.Println(err)
		log.Println(err)
		return -1

	}
	return db.UserNameToID(name)
}

func (db *DB) GetUserByName(username string) (User, error) {
	usr, err := db.query_row_UserSql(getUserByUsername, username)
	return usr, err
}

func (db *DB) UserNameToID(username string) int {
	usr, err := db.query_row_UserSql(getUserByUsername, username)
	if err != nil {
		return -1

	}
	return usr.ID
}
func (db *DB) GetUserByID(id int) (User, error) {
	return db.query_row_UserSql(getUserByID, id)
}

func (db *DB) GetUserGroupRelations(user_id int) ([]user_group_relations, error) {
	return db.query_rows_UserGroupRelationsSql(getUserGroupRelations, user_id)
}
func (db *DB) ListAllUserGroupRelations() ([]user_group_relations, error) {
	return db.query_rows_UserGroupRelationsSql(getAllUserGroupRelations)
}

// ListUsers fetches all users from the database.
func (db *DB) ListUsers() (map[string]User, error) {
	//query := "SELECT id, username, password_hash,  created_at FROM users"
	rows, err := db.SQL.QueryContext(DataBase.ctx, listUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	userMap := make(map[string]User)
	for rows.Next() {
		var u User
		var updatedAt string
		if err := rows.Scan(&u.ID, &u.Username, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		if u.CreatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
			return nil, fmt.Errorf("failed to parse timestamp for user %s: %w", u.Username, err)
		}
		userMap[u.Username] = u
	}

	return userMap, rows.Err()
}
