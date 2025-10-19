package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user account in the database.
type User struct {
	ID           int
	Username     string
	PasswordHash []byte
	Role         string
	UpdatedAt    string // formatted string for display
}

// AddUser hashes the password and inserts a new user record.
func (db *DB) AddUser(ctx context.Context, username string, password string, role string) error {
	if username == "" || password == "" {
		return errors.New("username and password cannot be empty")
	}

	hash, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now().Format(time.RFC3339)

	// Insert the new user
	_, err = db.Sql.ExecContext(ctx,
		"INSERT INTO users (username, password_hash, role, updated_at) VALUES (?, ?, ?, ?)",
		username, hash, role, now)

	if err != nil {
		// A database error (potentially unique constraint violation) occurred
		log.Printf("DB error adding user %s: %v", username, err)
		return errors.New("username already exists or database error occurred")
	}

	return nil
}

// UpdateUser updates an existing user's password and role.
// The password update is optional; if 'newPassword' is empty, only the role is updated.
func (db *DB) UpdateUser(ctx context.Context, userID int, newUsername, newPassword, newRole string) error {
	if newUsername == "" || newRole == "" {
		return errors.New("username and role cannot be empty")
	}

	now := time.Now().Format(time.RFC3339)
	var hash []byte
	var err error

	// If a new password is provided, generate a new hash.
	if newPassword != "" {
		hash, err = HashPassword(newPassword)
		if err != nil {
			return fmt.Errorf("failed to hash new password: %w", err)
		}
	}

	var result sql.Result
	if newPassword != "" {
		// Update username, password_hash, role, and updated_at
		result, err = db.Sql.ExecContext(ctx,
			"UPDATE users SET username=?, password_hash=?, role=?, updated_at=? WHERE id=?",
			newUsername, hash, newRole, now, userID)
	} else {
		// Update username, role, and updated_at (keep old password hash)
		result, err = db.Sql.ExecContext(ctx,
			"UPDATE users SET username=?, role=?, updated_at=? WHERE id=?",
			newUsername, newRole, now, userID)
	}

	if err != nil {
		log.Printf("DB error updating user %d: %v", userID, err)
		return errors.New("failed to update user (username may already exist)")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// DeleteUser removes a user record by ID.
func (db *DB) DeleteUser(ctx context.Context, userID int) error {
	result, err := db.Sql.ExecContext(ctx, "DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		log.Printf("DB error deleting user %d: %v", userID, err)
		return fmt.Errorf("database error during deletion: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// ListUsers fetches all users from the database.
// NOTE: The `currentUsername` parameter is not used here but is preserved
// to maintain the signature expected by HandleNewUser.
func (db *DB) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := db.Sql.QueryContext(ctx, "SELECT id, username, password_hash, role, updated_at FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return users, nil
}

// GetUserByID retrieves a single user by their ID.
func (db *DB) GetUserByID(ctx context.Context, userID int) (*User, error) {
	row := db.Sql.QueryRowContext(ctx, "SELECT id, username, password_hash, role, updated_at FROM users WHERE id = ?", userID)

	var u User
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to query user by ID: %w", err)
	}
	return &u, nil
}

// GetUserByUsername retrieves a single user by their username.
func (db *DB) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	row := db.Sql.QueryRowContext(ctx, "SELECT id, username, password_hash, role, updated_at FROM users WHERE username = ?", username)

	var u User
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to query user by username: %w", err)
	}
	return &u, nil
}

// HashPassword generates a bcrypt hash of the password.
func HashPassword(password string) ([]byte, error) {
	// Use a cost factor suitable for your environment. 10 is standard.
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// CheckPasswordHash compares a password with its hash.
func CheckPasswordHash(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}
