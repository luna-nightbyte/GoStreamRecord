package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// SQL QUERIES ---------------------------------------------------------------------

// AddUser hashes the password and inserts a new user record.
func (db *Tab) New(tabName, description string) error {
	if tabName == "" {
		return errors.New("tabName cannot be empty")
	}
	_, err := DataBase.SQL.ExecContext(DataBase.ctx, createTab, tabName, description)
	if err != nil {
		if strings.Contains(err.Error(), ErrIsExist) {
			return errors.New("tab already exists")
		}
		return err
	}

	return nil
}

// ListUsers fetches all users from the database.
func (db *Tab) List() (map[string]Tab, error) {
	//query := "SELECT id, username, password_hash,  created_at FROM users"
	rows, err := DataBase.SQL.QueryContext(DataBase.ctx, listTabs)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	tabMap := make(map[string]Tab)
	for rows.Next() {
		var u Tab
		if err := rows.Scan(&u.ID, &u.Name, &u.Description); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		tabMap[u.Name] = u
	}

	return tabMap, rows.Err()
}

// GetAvailableTabsForUser retrieves all tabs a user has access to.
// It takes a database connection pointer and InternalUsera user ID.
// This function replaces your original `GetAvalable` method.
func (db *Tab) GetAvailableTabsForUser(userID int) (map[string]Tab, error) {
	rows, err := DataBase.SQL.Query(getVisibleTabsForUser, userID)
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
func (db *Tab) ShareTab(tabID, groupID int) error {
	_, err := DataBase.SQL.ExecContext(DataBase.ctx, shareTabWithGroup, tabID, groupID)
	if err != nil {
		if strings.Contains(err.Error(), ErrIsExist) {
			return errors.New("Username exists")
		}
		return err
	}

	return nil
}

func (db *Tab) DeleteForGroup(groupID, tabID int) error {
	_, err := DataBase.SQL.ExecContext(DataBase.ctx, unshareTabFromGroup, tabID, groupID)
	return err
}

// HELPERS ------------------------------------------------------------------------------------
func (u *Tab) queryTabSql(query string, args ...any) error {
	row := DataBase.SQL.QueryRowContext(DataBase.ctx, query, args...)
	err := row.Scan(&u.ID, &u.Name, &u.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}

func (u *Tab) queryTabrGroupRelationsSql(query string, args ...any) (tab_group_relations, error) {
	row := DataBase.SQL.QueryRowContext(DataBase.ctx, query, args...)
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
