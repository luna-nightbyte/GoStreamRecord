package db

import (
	"database/sql"
	"errors"
	"fmt"
	"remoteCtrl/internal/utils"
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

	createdDate := time.Unix(int64(now.Unix()), 0)
	user_id := DataBase.Users.NameToID(username)
	expiresDate := createdDate.AddDate(0, 1, 0)

	_, err := DataBase.SQL.ExecContext(DataBase.ctx, createApi, apiName, user_id, utils.RandString(64), fmt.Sprint(expiresDate), fmt.Sprint(createdDate))
	if err != nil {
		if strings.Contains(err.Error(), ErrIsExist) {
			return errors.New("api already exists")
		}
		return err
	}
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
func (db *Api) List(owner_id int) (map[string]Api, error) {
	//query := "SELECT id, username, password_hash,  created_at FROM users"
	rows, err := DataBase.SQL.QueryContext(DataBase.ctx, listApis)
	if err != nil {
		return nil, fmt.Errorf("failed to query apis: %w", err)
	}
	defer rows.Close()

	tabMap := make(map[string]Api)
	for rows.Next() {
		var u Api
		if err := rows.Scan(&u.ID, &u.Name, &u.OwnerID, &u.Key, &u.Expires, &u.Created); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		if u.OwnerID != owner_id {
			continue
		}
		tabMap[u.Name] = u
	}

	return tabMap, rows.Err()
}

func (db *Api) DeleteForUser(user_id, api_id int) error {
	err := db.queryTabSql(deleteApi, user_id, api_id)
	return err
}

// HELPERS ------------------------------------------------------------------------------------
func (u *Api) queryTabSql(query string, args ...any) error {

	row := DataBase.SQL.QueryRowContext(DataBase.ctx, query, args...)
	err := row.Scan(&u.ID, &u.Name, &u.Key, &u.Expires, &u.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}
