package db

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	// "simserver/internal/util/env"
)

type DB struct {
	path string
	Sql  *sql.DB
}

type DataBaseTables struct {
	Users []UsersTable `json:"users"`
}

// UsersTable reflects the 'users' SQL table.
type UsersTable struct {
	// Note: We MUST use the string literal "id" here.
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash []byte    `json:"password_hash"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// User table column names
const (
	UserColID           = "id"
	UserColUsername     = "username"
	UserColPasswordHash = "password_hash"
	UserColRole         = "role"
	UserColUpdatedAt    = "updated_at"
)

var DataBase *DB
var AeadKey []byte

func (db *DB) Authenticate_(ctx context.Context, username, password string) (bool, error) {
	user, err := db.GetUserByUsername(ctx, username)
	if err != nil {
		return false, err
	}
	if CheckPasswordHash(password, user.PasswordHash) {
		return true, nil
	}

	return false, errors.New("invalid password")
}

// --- 2. SQL Schema Definition (Using the constant strings for clarity, though they are not required here) ---
func initSQLSchema() []string {
	// The SQL schema strings must be literal strings.
	return []string{
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS users(
            %s INTEGER PRIMARY KEY AUTOINCREMENT,
            %s TEXT NOT NULL UNIQUE,
            %s BLOB NOT NULL,
            %s TEXT NOT NULL,
            %s TEXT NOT NULL
        );`, UserColID, UserColUsername, UserColPasswordHash, UserColRole, UserColUpdatedAt),
	}
}

// Returns a mapping of Go struct field names (keys) to their corresponding
// database column names (values) using the defined constants.
func (u UsersTable) ColumnMap() map[string]string {
	return map[string]string{
		"ID":           UserColID,
		"Username":     UserColUsername,
		"PasswordHash": UserColPasswordHash,
		"Role":         UserColRole,
		"UpdatedAt":    UserColUpdatedAt,
	}
}
func Init(ctx context.Context, path string) {

	var db *DB
	db = new(DB)
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "simdb.sqlite"
	}
	if path != "" {
		dbPath = path
	}
	var err error
	db.path = dbPath
	db.Sql, err = Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.CreateDatabase(ctx); err != nil {
		log.Fatal(err)
	}
	if l, _ := db.ListUsers(ctx); len(l) == 0 {
		fmt.Println("Created default admin user: \nUsername: admin\nPassowrd: pass")
		if err := db.AddUser(ctx, "admin", "pass", "admin"); err != nil {
			log.Fatal("Cant add new user!")
		}
	}
	DataBase = db
}

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

func (db *DB) CreateDatabase(ctx context.Context) error {

	for _, s := range initSQLSchema() {
		if _, err := db.Sql.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func getSecret(secretKey string) []byte {

	secret, err := base64.StdEncoding.DecodeString(mustEnv(secretKey))
	if err != nil || len(secret) != 32 {
		log.Fatalf("SECRET_KEY must be base64 of 32 bytes")
	}
	return secret
}

func mustEnv(key string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		log.Fatalf("missing required env %s", key)
	}
	return v
}
