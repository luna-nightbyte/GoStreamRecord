package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

// ErrUserNotFound is returned when a user is not found in the database.
var ErrUserNotFound = errors.New("user not found")

// AddUser hashes the password and inserts a new user record.
func (db *DB) AddUser(username, password string, group string) error {
	if username == "" || password == "" {
		return errors.New("username and password cannot be empty")
	}

	hash, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now().Format(time.RFC3339)

	_, err = db.SQL.ExecContext(db.ctx, add_user, username, hash, now)
	if err != nil {
		log.Printf("DB error adding user %s: %v", username, err)
		return errors.New("username already exists or a database error occurred")
	}

	usrs, _ := db.ListUsers()

	_, err = db.SQL.ExecContext(db.ctx, admin_add_user, usrs[username].ID, group, now)
	if err != nil {
		log.Printf("DB error adding user %s: %v", username, err)
		return errors.New("username already exists or a database error occurred")
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
		hash, err := hashPassword(newPassword)
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
		log.Printf("DB error updating user %d: %v", userID, err)
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
	result, err := db.SQL.ExecContext(db.ctx, admin_del_user, userID)
	if err != nil {
		log.Printf("DB error deleting user %d: %v", userID, err)
		return fmt.Errorf("database error during deletion: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// ListUsers fetches all users from the database.
func (db *DB) ListUsers() (map[string]User, error) {
	//query := "SELECT id, username, password_hash,  created_at FROM users"
	rows, err := db.SQL.QueryContext(db.ctx, listUsers)
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

// GetUserByUsername retrieves a single user by their username.
func (db *DB) GetUserByUsername(username string) (*User, error) {
	row := db.SQL.QueryRowContext(db.ctx, getUserRoleInGroup, username)

	var u User
	var updatedAt string
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to query user by username: %w", err)
	}

	if u.CreatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
		return nil, fmt.Errorf("failed to parse timestamp for user %s: %w", u.Username, err)
	}
	return &u, nil
}

// Authenticate checks a user's credentials against the database.
func (db *DB) Authenticate(username, password string) (bool, error) {
	user, err := db.GetUserByUsername(username)
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

// IsAdmin checks if a user has admin privileges.
func (db *DB) IsAdmin(username string) (bool, error) {
	user, err := db.GetUserByUsername(username)
	if err != nil {
		return false, err
	}

	_, role, err := db.GetGroupForUser(strconv.Itoa(user.ID))
	if role == RoleAdmin {
		return true, nil
	}
	return false, err
}
