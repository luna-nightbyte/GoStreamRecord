package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// SQL QUERIES ---------------------------------------------------------------------

// AddUser hashes the password and inserts a new user record.
func (db *Api) New(apiName, username string) error {
	if apiName == "" {
		return errors.New("tabName cannot be empty")
	}
	now := time.Now()
	user_id := DataBase.Users.NameToID(username)
	_, err := DataBase.SQL.ExecContext(DataBase.ctx, createApi, apiName, "randomKEy", now.Unix(), now.Unix())
	if err != nil {
		if strings.Contains(err.Error(), ErrIsExist) {
			return errors.New("tab already exists")
		}
		return err
	}
	apis, _ := db.List()
	db.share(apis[apiName].ID, user_id)
	return nil
}

// ListUsers fetches all users from the database.
func (db *Api) ListUserApis(user_id int) (map[string]Api, error) {
	//query := "SELECT id, username, password_hash,  created_at FROM users"
	rows, err := DataBase.SQL.QueryContext(DataBase.ctx, getUserApis, user_id)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	tabMap := make(map[string]Api)
	for rows.Next() {
		var u Api
		if err := rows.Scan(&u.ID, &u.Name, &u.Key, &u.Expires, &u.Created); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		tabMap[u.Name] = u
	}

	return tabMap, rows.Err()
}

// ListUsers fetches all users from the database.
func (db *Api) List() (map[string]Api, error) {
	//query := "SELECT id, username, password_hash,  created_at FROM users"
	rows, err := DataBase.SQL.QueryContext(DataBase.ctx, listApis)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	tabMap := make(map[string]Api)
	for rows.Next() {
		var u Api
		if err := rows.Scan(&u.ID, &u.Name, &u.Key, &u.Expires, &u.Created); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		tabMap[u.Name] = u
	}

	return tabMap, rows.Err()
}

// GetAvailableTabsForUser retrieves all tabs a user has access to.
// It takes a database connection pointer and InternalUsera user ID.
// This function replaces your original `GetAvalable` method.
func (db *Api) GetAvailableTabsForUser(userID int) (map[string]Tab, error) {
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
func (db *Api) DeleteTabForUser(user_id, tab_id int) (*Api, error) {
	err := db.queryTabSql(unshareTabFromGroup, user_id, tab_id)
	return db, err
}

// HELPERS ------------------------------------------------------------------------------------
func (u *Api) queryTabSql(query string, args ...any) error {

	var created, expires string
	row := DataBase.SQL.QueryRowContext(DataBase.ctx, query, args...)
	err := row.Scan(&u.ID, &u.Name, &u.Key, &expires, &created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	if u.Created, err = time.Parse(time.RFC3339, created); err != nil {
		return err
	}
	if u.Expires, err = time.Parse(time.RFC3339, expires); err != nil {
		return err
	}

	return nil
}

func (u *Api) queryTabrGroupRelationsSql(query string, args ...any) (user_api_relations, error) {
	row := DataBase.SQL.QueryRowContext(DataBase.ctx, query, args...)
	var usrGrp user_api_relations
	err := row.Scan(&usrGrp.ApiID, &usrGrp.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return usrGrp, ErrUserNotFound
		}
		return usrGrp, err
	}

	return usrGrp, nil
}

// AddGroup inserts a new group with a given set of permissions.
func (db *Api) share(apiID, userID int) error {
	_, err := DataBase.SQL.ExecContext(DataBase.ctx, createApiRelation, userID, apiID)
	if err != nil {
		if strings.Contains(err.Error(), ErrIsExist) {
			return errors.New("Username exists")
		}
		return err
	}

	return nil
}
