package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// SQL QUERIES ---------------------------------------------------------------------

func (db *DB) NewStreamer(streamer_name, provider string, user_id int, share bool) {
	db.newStreamer(streamer_name, provider, user_id)
	if share {
		groups, _, _ := db.ListGroupsByUserID(user_id)
		streamers, _ := db.ListStreamers()
		for _, group := range groups {
			db.ShareStreamer(streamers[streamer_name].ID, group.ID)
		}
	}
}
// AddUser hashes the password and inserts a new user record.
func (db *DB) newStreamer(streamerName, provider string, user_id int) error {
	if streamerName == "" {
		return errors.New("name cannot be empty")
	}
	_, err := db.SQL.ExecContext(db.ctx, createStreamer, streamerName, provider, user_id)
	if err != nil {
		if strings.Contains(err.Error(), ErrIsExist) {
			return errors.New("already exists")
		}
		return err
	}

	return nil
}

// GetAvailableStreamersForUser retrieves all tabs a user has access to.
// It takes a database connection pointer and InternalUsera user ID. 
func (db *DB) GetAvailableStreamersForUser(userID int) (map[string]Streamer, error) {
	rows, err := db.SQL.Query(getVisibleStreamerForUser, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	streamersMap := make(map[string]Streamer)
	for rows.Next() {
		var s Streamer
		if err := rows.Scan(&s.ID, &s.Name, &s.Provider, &s.UploaderUserID); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		streamersMap[s.Name] = s
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return streamersMap, nil
}
func (db *DB) GetAvailableStreamersForGroup(groupID int) (map[string]Streamer, error) {
	rows, err := db.SQL.Query(getVisibleStreamerForGroup, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	tabsMap := make(map[string]Streamer)
	for rows.Next() {
		var s Streamer
		if err := rows.Scan(&s.ID, &s.Name, &s.Provider, &s.UploaderUserID); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		tabsMap[s.Name] = s
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return tabsMap, nil
}
func (db *DB) DeleteStreamerForUser(user_id, streamer_id int) ( error) { 
	_, err := db.SQL.ExecContext(db.ctx, removeUploaderUserFromStreamer, streamer_id, user_id)
	return  err
}

func (db *DB) DeleteStreamerForGroup(groupID, streamerID int) error {
	_, err := db.SQL.ExecContext(db.ctx, unshareStreamerFromGroup, streamerID, groupID)
	return err
}

// AddGroup inserts a new group with a given set of permissions.
func (db *DB) ShareStreamer(streamerID, groupID int) error {
	_, err := db.SQL.ExecContext(db.ctx, shareStreamerWithGroup, streamerID, groupID)
	if err != nil {
		if strings.Contains(err.Error(), ErrIsExist) {
			return errors.New("Username exists")
		}
		return err
	}

	return nil
}

// ListUsers fetches all users from the db.
func (db *DB) ListStreamers() (map[string]Streamer, error) {
	//query := "SELECT id, username, password_hash,  created_at FROM users"
	rows, err := db.SQL.QueryContext(db.ctx, listStreamer)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	streamerMap := make(map[string]Streamer)
	for rows.Next() {
		var s Streamer
		if err := rows.Scan(&s.ID, &s.Name, &s.Provider, &s.UploaderUserID); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		streamerMap[s.Name] = s
	}

	return streamerMap, rows.Err()
}

// HELPERS ------------------------------------------------------------------------------------
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
