package db

import (
	"errors"
	"fmt"
	"strings"
)

// SQL QUERIES ---------------------------------------------------------------------

// AddUser hashes the password and inserts a new user record.
func (db *DB) NewTab(tabName, description string) error {
	if tabName == "" {
		return errors.New("tabName cannot be empty")
	}
	_, err := db.SQL.ExecContext(db.ctx, createTab, tabName, description)
	if err != nil {
		if strings.Contains(err.Error(), ErrIsExist) {
			return errors.New("tab already exists")
		}
		return err
	}

	return nil
}

// ListUsers fetches all users from the db.
func (db *DB) ListTabs() (map[string]Tab, error) {
	//query := "SELECT id, username, password_hash,  created_at FROM users"
	rows, err := db.SQL.QueryContext(db.ctx, listTabs)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	tabMap := make(map[string]Tab)
	for rows.Next() {
		var t Tab
		if err := rows.Scan(&t.ID, &t.Name, &t.Description); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		tabMap[t.Name] = t
	}

	return tabMap, rows.Err()
}

// GetAvailableTabsForUser retrieves all tabs a user has access to.
// It takes a database connection pointer and InternalUsera user ID.
// This function replaces your original `GetAvalable` method.
func (db *DB) GetAvailableTabsForUser(userID int) (map[string]Tab, error) {
	rows, err := db.SQL.Query(getVisibleTabsForUser, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getVisibleTabsForUser query: %w", err)
	}
	defer rows.Close()
	tabsMap := make(map[string]Tab)
	for rows.Next() {
		var t Tab
		if err := rows.Scan(&t.ID, &t.Name, &t.Description); err != nil {
			return nil, fmt.Errorf("failed to scan tab row: %w", err)
		}
		tabsMap[t.Name] = t
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during tab row iteration: %w", err)
	}

	return tabsMap, nil
}

// AddGroup inserts a new group with a given set of permissions.
func (db *DB) ShareTab(tabID, groupID int) error {
	_, err := db.SQL.ExecContext(db.ctx, shareTabWithGroup, tabID, groupID)
	if err != nil {
		if strings.Contains(err.Error(), ErrIsExist) {
			return errors.New("Username exists")
		}
		return err
	}

	return nil
}

func (db *DB) DeleteTabForGroup(groupID, tabID int) error {
	_, err := db.SQL.ExecContext(db.ctx, unshareTabFromGroup, tabID, groupID)
	return err
}
