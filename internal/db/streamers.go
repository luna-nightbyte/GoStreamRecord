package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// SQL QUERIES ---------------------------------------------------------------------

// AddUser hashes the password and inserts a new user record.
func (db *Streamer) New(tabName, provider string) error {
	if tabName == "" {
		return errors.New("tabName cannot be empty")
	}
	_, err := DataBase.SQL.ExecContext(DataBase.ctx, createStreamer, tabName, provider)
	if err != nil {
		if strings.Contains(err.Error(), ErrIsExist) {
			return errors.New("tab already exists")
		}
		return err
	}

	return nil
}

// ListUsers fetches all users from the database.
func (db *Streamer) List() (map[string]Streamer, error) {
	//query := "SELECT id, username, password_hash,  created_at FROM users"
	rows, err := DataBase.SQL.QueryContext(DataBase.ctx, listStreamer)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	tabMap := make(map[string]Streamer)
	for rows.Next() {
		var u Streamer
		if err := rows.Scan(&u.ID, &u.Name, &u.Provider); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		tabMap[u.Name] = u
	}

	return tabMap, rows.Err()
}

// GetAvailableTabsForUser retrieves all tabs a user has access to.
// It takes a database connection pointer and InternalUsera user ID.
// This function replaces your original `GetAvalable` method.
func (db *Streamer) GetAvailableForUser(userID int) (map[string]Streamer, error) {
	rows, err := DataBase.SQL.Query(getVisibleStreamerForUser, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getVisibleTabsForUser query: %w", err)
	}
	defer rows.Close()
	tabsMap := make(map[string]Streamer)
	for rows.Next() {
		var t Streamer
		if err := rows.Scan(&t.ID, &t.Name, &t.Provider); err != nil {
			return nil, fmt.Errorf("failed to scan tab row: %w", err)
		}
		tabsMap[t.Name] = t
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during tab row iteration: %w", err)
	}

	return tabsMap, nil
}
func (db *Streamer) DeleteForUser(user_id, tab_id int) (*Streamer, error) {
	err := db.queryTabSql(unshareStreamerFromGroup, user_id, tab_id)
	return db, err
}

// HELPERS ------------------------------------------------------------------------------------
func (u *Streamer) queryTabSql(query string, args ...any) error {
	row := DataBase.SQL.QueryRowContext(DataBase.ctx, query, args...)
	err := row.Scan(&u.ID, &u.Name, &u.Provider)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}

func (u *Streamer) queryTabrGroupRelationsSql(query string, args ...any) (streamer_group_relations, error) {
	row := DataBase.SQL.QueryRowContext(DataBase.ctx, query, args...)
	var usrGrp streamer_group_relations
	err := row.Scan(&usrGrp.StreamerID, &usrGrp.GroupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return usrGrp, ErrUserNotFound
		}
		return usrGrp, err
	}

	return usrGrp, nil
}

// AddGroup inserts a new group with a given set of permissions.
func (db *Streamer) Share(streamerID, groupID int) error {
	_, err := DataBase.SQL.ExecContext(DataBase.ctx, shareStreamerbWithGroup, streamerID, groupID)
	if err != nil {
		if strings.Contains(err.Error(), ErrIsExist) {
			return errors.New("Username exists")
		}
		return err
	}

	return nil
}
