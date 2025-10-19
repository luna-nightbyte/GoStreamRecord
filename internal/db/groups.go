package db

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// AddGroup inserts a new group with a given set of permissions.
func (db *DB) AddGroup(groupName string, description string) error {
	if groupName == "" {
		return errors.New("group name cannot be empty")
	}

	now := time.Now().Format(time.RFC3339)

	//query := "INSERT INTO groups (name, permissions, updated_at) VALUES (?, ?, ?)"
	_, err := db.SQL.ExecContext(db.ctx, createGroup, groupName, description, now)
	if err != nil {
		log.Printf("DB error adding group %s: %v", groupName, err)
		return errors.New("group name already exists or a database error occurred")
	}

	return nil
}

// ListGroups fetches all groups from the database.
func (db *DB) ListGroups() (map[string]Group, error) {
	// query := "SELECT id, name, permissions, updated_at FROM groups"
	rows, err := db.SQL.QueryContext(db.ctx, getGroupByGroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to query groups: %w", err)
	}
	defer rows.Close()

	groupMap := make(map[string]Group)
	for rows.Next() {
		var g Group
		var updatedAt string
		if err := rows.Scan(&g.ID, &g.Name, &g.Description, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan group row: %w", err)
		}

		if g.CreatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
			return nil, fmt.Errorf("failed to parse timestamp for group %s: %w", g.Name, err)
		}
		groupMap[g.Name] = g
	}

	return groupMap, rows.Err()
}

func (db *DB) GetGroupForUser(user_id string) (map[string]Group, string, error) {
	rows, err := db.SQL.QueryContext(db.ctx, getGroupsForUser, user_id)
	if err != nil {
		return nil, "", fmt.Errorf("failed to query groups: %w", err)
	}
	defer rows.Close()

	var role string
	groupMap := make(map[string]Group)
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.Name, &g.Description, &role); err != nil {
			return nil, "", fmt.Errorf("failed to scan group row: %w", err)
		}

		groupMap[g.Name] = g
	}

	return groupMap, role, rows.Err()
}
