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
func (g *Group) New(groupName string, description string) error {
	if groupName == "" {
		return errors.New("group name cannot be empty")
	}

	now := time.Now().Format(time.RFC3339)

	//query := "INSERT INTO groups (name, permissions, updated_at) VALUES (?, ?, ?)"
	_, err := DataBase.SQL.ExecContext(DataBase.ctx, createGroup, groupName, description, now)
	if err != nil {
		log.Printf("DB error adding group %s: %v", groupName, err)
		return errors.New("group name already exists or a database error occurred")
	}

	return nil
}

// AddGroup inserts a new group with a given set of permissions.
func (g *Group) AddUser(userID, groupID int, role string) error {

	_, err := DataBase.SQL.ExecContext(DataBase.ctx, addUserToGroup, userID, groupID, role)
	if err != nil {
		if strings.Contains(err.Error(), ErrIsExist) {
			return errors.New("Username exists")
		}
		return err
	}

	return nil
}

// AddGroup inserts a new group with a given set of permissions.
func (g *Group) RemoveUser(userID, groupID int) error {

	_, err := DataBase.SQL.ExecContext(DataBase.ctx, removeUserFromGroup, userID, groupID)
	if err != nil {
		if strings.Contains(err.Error(), ErrIsExist) {
			return errors.New("Username exists")
		}
		return err
	}

	return nil
}

// ListGroups fetches all groups from the database.
func (g *Group) List() (map[string]Group, error) {
	// query := "SELECT id, name, permissions, updated_at FROM groups"
	rows, err := DataBase.SQL.QueryContext(DataBase.ctx, getAllGroups)
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

func (g *Group) ListGroupsByUserID(user_id int) (map[string]Group, string, error) {
	rows, err := DataBase.SQL.QueryContext(DataBase.ctx, getGroupsForUser, user_id)
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

// GetUserByUsername retrieves a single user by their username.
func (u *Group) GetGroupByName(username string) (*Group, error) {
	err := u.queryGroupSql(getGroupByName, username)
	return u, err
}

func (g *Group) NameToID(groupName string) int {
	grps, err := g.List()
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
func (g *Group) queryGroupSql(query string, args ...any) error {
	row := DataBase.SQL.QueryRowContext(DataBase.ctx, query, args...)

	var createdAt string
	err := row.Scan(&g.ID, &g.Name, &g.Description, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	if g.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
		return err
	}
	return nil
}
