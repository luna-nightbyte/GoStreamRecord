package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"remoteCtrl/internal/utils"
)

// DB wraps the sql.DB connection pool.
type DB struct {
	ctx    context.Context
	SQL    *sql.DB
	Users  User
	Groups Group
	V      Video
}

// Global variable to hold the database instance.
var DataBase *DB

// Internal server user
const InternalUser string = "_internal"

// createSchema executes the necessary SQL statements to build the database tables.
func createSchema(ctx context.Context, db *sql.DB) error {
	schemaSQL := fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
		q_create_users,
		q_create_goups,
		q_create_user_group_roles,
		q_create_videos,
		q_create_video_groups,
	)

	if _, err := db.ExecContext(ctx, schemaSQL); err != nil {
		return fmt.Errorf("failed to execute schema creation: %w", err)
	}
	return nil
}

// Init initializes the database connection, creates the schema if it doesn't exist,
// and ensures a default admin user is present.
func Init(ctx context.Context, path string) {

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "database.sqlite"
	}
	// Allow overriding path for testing or specific configurations.
	if path != "" {
		dbPath = path
	}

	_, err := os.Stat(dbPath)
	isNewDb := os.IsNotExist(err)

	db, err := open(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	DataBase = &DB{
		ctx: ctx,
		SQL: db}
	if isNewDb {
		// Create database
		if err := createSchema(ctx, DataBase.SQL); err != nil {
			log.Fatalf("Failed to create database schema: %v", err)
		}
		// Groups

		if err := DataBase.Groups.New(GroupDefault, "admins with full control"); err != nil {
			log.Fatalf("Fatal: Could not create default group: %v", err)
		}
		// Users
		var defaultUser, defaultPass string = "user", "password"
		if err := DataBase.Users.New(defaultUser, defaultPass); err != nil {
			log.Fatalf("Fatal: Could not add default admin user: %v", err)
		}
		// Add to group
		user_id := DataBase.Users.NameToID(defaultUser)
		admin_group_id := DataBase.Groups.NameToID(GroupDefault)

		DataBase.Groups.AddUser(user_id, admin_group_id, RoleAdmin)

		if err := DataBase.Groups.AddUser(user_id, admin_group_id, RoleAdmin); err != nil {
			log.Fatalf("Fatal: Could not add default admin user: %v", err)
		}

		randPass, _ := hashPassword(utils.RandString(15))
		DataBase.Users.New(InternalUser, string(randPass))

		user_id = DataBase.Users.NameToID(InternalUser)

		DataBase.Groups.AddUser(user_id, admin_group_id, RoleAdmin)

		fmt.Println("New database created. Creating default admin user 'user' with password 'password'")
	}
}

// Open establishes a new database connection.
func open(path string) (*sql.DB, error) {
	// Using "_pragma=foreign_keys(1)" to enforce foreign key constraints.
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_pragma=foreign_keys(1)", path))
	if err != nil {
		return nil, fmt.Errorf("could not open sqlite database at %s: %w", path, err)
	}

	// Set a single connection to avoid race conditions with file-based SQLite.
	db.SetMaxOpenConns(1)

	return db, nil
}
