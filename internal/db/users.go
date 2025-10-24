package db

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"remoteCtrl/internal/web/handlers/cookie"
	"strings"
	"time"
)

// ErrUserNotFound is returned when a user is not found in the database.
var ErrUserNotFound = errors.New("user not found")

const ErrIsExist = "UNIQUE constraint failed"

// SQL QUERIES ---------------------------------------------------------------------

// AddUser hashes the password and inserts a new user record.
func (db *User) New(username, raw_password string) error {
	if username == "" || raw_password == "" {
		return errors.New("username and password cannot be empty")
	}
	hash, err := hashPassword(raw_password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	now := time.Now().Format(time.RFC3339)
	_, err = DataBase.SQL.ExecContext(DataBase.ctx, createUser, username, hash, now)
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
func (db *User) Update(userID int, newUsername string, newPassword string) error {
	if newUsername == "" {
		return errors.New("username cannot be empty")
	}

	var result sql.Result
	var err error
	if newPassword != "" {
		hash, err := hashPassword(newPassword)
		if err != nil {
			return fmt.Errorf("failed to hash new password: %w", err)
		}
		result, err = DataBase.SQL.ExecContext(DataBase.ctx, updateUser, newUsername, hash, userID)
	} else {
		usrs, _ := DataBase.Users.List()
		for _, urs := range usrs {
			if urs.ID == userID {
				result, err = DataBase.SQL.ExecContext(DataBase.ctx, updateUser, newUsername, urs.PasswordHash, userID)

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
func (db *User) Delete(userID int) error {
	result, err := DataBase.SQL.ExecContext(DataBase.ctx, admin_del_user, userID)
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
func (db *User) Authenticate(username, password string) (bool, error) {
	user, err := db.GetUserByName(username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return false, errors.New("invalid username or password")
		}
		return false, err
	}

	if checkPasswordHash(password, user.PasswordHash) {
		return true, nil
	}

	return false, errors.New("invalid username or password")
}

// IsAdmin checks ifunc (db *User) Authenticate(username, password string) (bool, error)f a user has admin privileges.
func (db *User) IsAdmin(username string) (bool, error) {
	user, err := db.GetUserByName(username)
	if err != nil {
		return false, err
	}
	_, role, err := DataBase.Groups.ListGroupsByUserID(user.ID)
	if role == RoleAdmin {
		return true, nil
	}
	return false, err
}

func (db *User) HttpRequestID(r *http.Request) int {
	name, _ := cookie.ValidateSession(r)
	return db.NameToID(name)
}

func (db *User) GetUserByName(username string) (*User, error) {
	err := db.queryUserSql(getUserByUsername, username)
	return db, err
}

func (u *User) NameToID(username string) int {
	u.queryUserSql(getUserByUsername, username)
	return u.ID
}
func (db *User) GetUserByID(id int) (*User, error) {
	err := db.queryUserSql(getUserByID, id)
	return db, err
}

// ListUsers fetches all users from the database.
func (db *User) List() (map[string]User, error) {
	//query := "SELECT id, username, password_hash,  created_at FROM users"
	rows, err := DataBase.SQL.QueryContext(DataBase.ctx, listUsers)
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

// HELPERS ------------------------------------------------------------------------------------
func (u *User) queryUserSql(query string, args ...any) error {
	row := DataBase.SQL.QueryRowContext(DataBase.ctx, query, args...)

	var createdAt string
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	if u.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
		return err
	}
	return nil
}

func (u *User) queryUserGroupRelationsSql(query string, args ...any) (user_group_relations, error) {
	row := DataBase.SQL.QueryRowContext(DataBase.ctx, query, args...)
	var usrGrp user_group_relations
	err := row.Scan(&usrGrp.UserID, &usrGrp.GroupID, &usrGrp.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return usrGrp, ErrUserNotFound
		}
		return usrGrp, err
	}

	return usrGrp, nil
}
