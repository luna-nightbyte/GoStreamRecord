package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

//  GROUPS ----------------------------------------------------------------------------------------

// AddGroup inserts a new group with a given set of permissions.
func (db *DB) NewGroup(groupName string, description string) error {
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

// AddGroup inserts a new group with a given set of permissions.
func (db *DB) AddUserToGroup(userID, groupID int, role string) error {

	_, err := db.SQL.ExecContext(db.ctx, addUserToGroup, userID, groupID, role)
	if err != nil {
		if strings.Contains(err.Error(), ErrIsExist) {
			return errors.New("Username exists")
		}
		return err
	}

	return nil
}

// AddGroup inserts a new group with a given set of permissions.
func (db *DB) RemoveUserFromGroup(userID, groupID int) error {

	_, err := db.SQL.ExecContext(db.ctx, removeUserFromGroup, userID, groupID)
	if err != nil {
		if strings.Contains(err.Error(), ErrIsExist) {
			return errors.New("Username exists")
		}
		return err
	}

	return nil
}

// ListGroups fetches all groups from the db.
func (db *DB) ListAllGroups() (map[string]Group, error) {
	// query := "SELECT id, name, permissions, updated_at FROM groups"
	rows, err := db.SQL.QueryContext(db.ctx, getAllGroups)
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

func (db *DB) ListGroupsByUserID(user_id int) (map[int]Group, string, error) {
	rows, err := db.SQL.QueryContext(db.ctx, getGroupsForUser, user_id)
	if err != nil {
		return nil, "", fmt.Errorf("failed to query groups: %w", err)
	}
	defer rows.Close()

	var role string
	groupMap := make(map[int]Group)
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.Name, &g.Description, &role); err != nil {
			return nil, "", fmt.Errorf("failed to scan group row: %w", err)
		}

		groupMap[g.ID] = g
	}

	return groupMap, role, rows.Err()
}

// GetUserByUsername retrieves a single user by their username.
func (db *DB) GetGroupByName(username string) (Group, error) {
	return db.queryGroupSql(getGroupByName, username)
}

// GetUserByUsername retrieves a single user by their username.
func (db *DB) GetGroupByID(id int) (Group, error) {
	return db.queryGroupSql(getGroupByGroupID, id)
}

func (db *DB) GroupNameToID(groupName string) int {
	grps, err := db.ListAllGroups()
	if err != nil {
		return -1
	}
	group, exists := grps[groupName]
	if !exists {
		return -1
	}
	return group.ID
}

// // GetUserByUsername retrieves a single user by their username.
// func (db *User) GetAllUserGroupRelations() (user_group_relations, error) {
// 	return db.queryUserGroupRelationsSql(getUserGroupRelations)
// }

// HELPERS ------------------------------------------------------------------------------------
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
