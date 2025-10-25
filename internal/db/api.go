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
func (db *DB) NewApi(apiName, username string) error {
	if apiName == "" {
		return errors.New("tabName cannot be empty")
	}
	now := time.Now()

	createdDate := time.Unix(int64(now.Unix()), 0)
	user_id := db.UserNameToID(username)
	expiresDate := createdDate.AddDate(0, 1, 0)

	err := db.execQuery(createApi, apiName, user_id, utils.RandString(64), fmt.Sprint(expiresDate), fmt.Sprint(createdDate))
	if err != nil {
		if strings.Contains(err.Error(), ErrIsExist) {
			return errors.New("api already exists")
		}
		return err
	}
	return nil
}

// ListUsers fetches all users from the db.
func (db *DB) ListAvailableAPIsForUser(user_id int) (map[string]Api, error) {
	//query := "SELECT id, username, password_hash,  created_at FROM users"
	rows, err := db.query(getUserApis, user_id)
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

func (db *DB) DeleteApiForUser(user_id, api_id int) error {
	return db.execQuery(deleteApi, user_id, api_id)
}

// HELPERS ------------------------------------------------------------------------------------
func (db *DB) queryApiSql(query string, args ...any) (Api, error) {
	var a Api
	row := db.SQL.QueryRowContext(db.ctx, query, args...)
	err := row.Scan(&a.ID, &a.Name, &a.Key, &a.Expires, &a.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return a, ErrNotFound
		}
		return a, err
	}

	return a, nil
}
