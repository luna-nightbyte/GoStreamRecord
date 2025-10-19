// Package db handles all database interactions for the application.
package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DB wraps the sql.DB connection pool.
type DB struct {
	ctx context.Context
	SQL *sql.DB
}

// Global variable to hold the database instance.
var DataBase *DB

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
	db, err := Open(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if err := createSchema(ctx, db); err != nil {
		log.Fatalf("Failed to create database schema: %v", err)
	}
	DataBase = &DB{
		ctx: ctx,
		SQL: db}

	if isNewDb {
		fmt.Println("New database created. Creating default admin user 'admin' with password 'password'")
		if err := DataBase.AddUser("admin", "password", RoleAdmin); err != nil {
			log.Fatalf("Fatal: Could not add default admin user: %v", err)
		}
		if err := DataBase.AddGroup(RoleAdmin, ""); err != nil {
			log.Fatalf("Fatal: Could not create default group: %v", err)
		}
		if err := DataBase.AddGroup(RoleViewers, ""); err != nil {
			log.Fatalf("Fatal: Could not create default group: %v", err)
		}
		if err := DataBase.AddGroup(RoleUsers, ""); err != nil {
			log.Fatalf("Fatal: Could not create default group: %v", err)
		}
		if err := DataBase.AddGroup(RoleRestricted, ""); err != nil {
			log.Fatalf("Fatal: Could not create default group: %v", err)
		}
	}
	// Check if any users exist. If not, create a default admin.
	users, err := DataBase.ListUsers()
	if err != nil {
		log.Fatalf("Failed to check for existing users: %v", err)
	}

	if len(users) == 0 {
	}
}

// Open establishes a new database connection.
func Open(path string) (*sql.DB, error) {
	// Using "_pragma=foreign_keys(1)" to enforce foreign key constraints.
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_pragma=foreign_keys(1)", path))
	if err != nil {
		return nil, fmt.Errorf("could not open sqlite database at %s: %w", path, err)
	}

	// Set a single connection to avoid race conditions with file-based SQLite.
	db.SetMaxOpenConns(1)

	return db, nil
}

// createSchema executes the necessary SQL statements to build the database tables.
func createSchema(ctx context.Context, db *sql.DB) error {
	schemaSQL := fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
		create_users,
		create_goups,
		create_user_group_roles,
		create_videos,
		create_video_groups,
	)

	if _, err := db.ExecContext(ctx, schemaSQL); err != nil {
		return fmt.Errorf("failed to execute schema creation: %w", err)
	}
	return nil
}

// Models & Table Structures

type User struct {
	ID           int       `json:"id"`
	PasswordHash []byte    `json:"-"` // Omit from JSON responses
	Username     string    `json:"username"`
	CreatedAt    time.Time `json:"created_at"`
}

type Group struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Video struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Filepath       string    `json:"filepath"`
	UploaderUserID int       `json:"uploader_user_id"`
	CreatedAt      time.Time `json:"created_at"`
}

// Roles define broad levels of access.
const (
	RoleAdmin      string = "admin"
	RoleUsers      string = "users"
	RoleViewers    string = "viewers"
	RoleRestricted string = "restricted"
)

// Permissions define granular actions.
const (
	PermNone   string = "-1"
	PermAll    string = "1"
	PermWrite  string = "2"
	PermRead   string = "3"
	PermDelete string = "4"
)

var (
	groupRestricted = []string{PermNone}
	groupReadOnly   = []string{PermRead}
	groupReadWrite  = []string{PermRead, PermWrite}
	groupAdmins     = []string{PermAll, PermRead, PermWrite, PermDelete}

	// GROUP maps a role constant to its set of permissions.
	GROUP = map[string][]string{
		RoleAdmin:      groupAdmins,
		RoleUsers:      groupReadWrite,
		RoleViewers:    groupReadOnly,
		RoleRestricted: groupRestricted,
	}
)

// initial queries
const (
	create_users string = `CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash BLOB NOT NULL,
    created_at TEXT NOT NULL
);`
	create_goups string = `CREATE TABLE groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TEXT NOT NULL
);`
	create_videos string = `CREATE TABLE videos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    filepath TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    uploader_user_id INTEGER NOT NULL,
    created_at TEXT NOT NULL,
    FOREIGN KEY (uploader_user_id) REFERENCES users (id)
);`
	create_user_group_roles string = `CREATE TABLE user_group_roles (
    user_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    role TEXT NOT NULL,  
    PRIMARY KEY (user_id, group_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (group_id) REFERENCES groups (id)
);`
	create_video_groups string = `CREATE TABLE video_groups (
    video_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    PRIMARY KEY (video_id, group_id),
    FOREIGN KEY (video_id) REFERENCES videos (id),
    FOREIGN KEY (group_id) REFERENCES groups (id)
);`
)

// Reusable queries
const (
	select_shared_videos string = `
	-- Select videos shared with the user's groups
	SELECT v.id, v.name, v.filepath
	FROM videos v
	JOIN video_groups vg ON v.id = vg.video_id
	JOIN user_group_roles ugr ON vg.group_id = ugr.group_id
	WHERE ugr.user_id = ? -- The logged-in user's ID

	UNION -- Combine results and remove duplicates

	-- Select videos uploaded by the user
	SELECT v.id, v.name, v.filepath
	FROM videos v
	WHERE v.uploader_user_id = ?; -- The logged-in user's ID`

	admin_del_user string = `DELETE FROM users WHERE user_id = ?`
	get_users      string = `SELECT user_id, group_id, role FROM user_group_roles WHERE user_id = ?`

	add_user string = `INSERT INTO users (username, password_hash, created_at) VALUES (?, ?, ?)`
	// user_id, group_id, role
	admin_add_user string = `INSERT INTO user_group_roles (user_id, group_id, role) VALUES (?, ?, ?);`

	mod_share_video string = `INSERT INTO video_groups (video_id, group_id) VALUES (?, ?);` // video_id, group_id
)

// This file centralizes all SQL queries for maintainability.
// Using constants for queries makes the Go code cleaner and prevents typos.

const (
	// --- User Queries (users table) ---

	// createUser inserts a new user record.
	createUser = `INSERT INTO users (username, password_hash, created_at) VALUES (?, ?, ?)`

	// getUserByID retrieves a single user with their password hash for authentication.
	getUserByID = `SELECT id, username, password_hash, created_at FROM users WHERE id = ?`

	// getUserByUsername retrieves a single user with their password hash for authentication.
	getUserByUsername = `SELECT id, username, password_hash, created_at FROM users WHERE username = ?`

	// listUsers retrieves all users without their password hashes for general listings.
	listUsers = `SELECT id, username, created_at FROM users ORDER BY username`

	// updateUser can change a user's username and password hash.
	updateUser = `UPDATE users SET username = ?, password_hash = ? WHERE id = ?`

	// deleteUser removes a user. Cascading deletes will handle their relationships.
	deleteUser = `DELETE FROM users WHERE id = ?`

	// --- Group Queries (groups table) ---

	// createGroup inserts a new group.
	createGroup = `INSERT INTO groups (name, description, created_at) VALUES (?, ?, ?)`

	// getGroupByID retrieves a single group by its primary key.
	getGroupByGroupID = `SELECT id, name, description, created_at FROM groups WHERE id = ?`

	// getGroupByName retrieves a single group by its unique name.
	getGroupByName = `SELECT id, name, description, created_at FROM groups WHERE name = ?`

	// listGroups retrieves all groups.
	listGroups = `SELECT id, name, description, created_at FROM groups ORDER BY name`

	// updateGroup changes a group's name or description.
	updateGroup = `UPDATE groups SET name = ?, description = ? WHERE id = ?`

	// deleteGroup removes a group. Cascading deletes will handle its relationships.
	deleteGroup = `DELETE FROM groups WHERE id = ?`

	// --- Video Queries (videos table) ---

	// createVideo inserts metadata for a new video.
	createVideo = `INSERT INTO videos (filepath, name, uploader_user_id, created_at) VALUES (?, ?, ?, ?)`

	// getVideoByID retrieves a single video by its ID.
	getVideoByID = `SELECT id, filepath, name, uploader_user_id, created_at FROM videos WHERE id = ?`

	// listAllVideos retrieves all videos, useful for an admin view.
	listAllVideos = `SELECT id, filepath, name, uploader_user_id, created_at FROM videos ORDER BY name`

	// deleteVideo removes a video. Cascading deletes will handle its shares.
	deleteVideo = `DELETE FROM videos WHERE id = ?`

	// --- Relationship Queries (user_group_roles table) ---

	// addUserToGroup assigns a user to a group with a specific role.
	// 'INSERT OR REPLACE' is used to easily change a user's role.
	addUserToGroup = `INSERT OR REPLACE INTO user_group_roles (user_id, group_id, role) VALUES (?, ?, ?)`

	// removeUserFromGroup removes a user's membership from a group.
	removeUserFromGroup = `DELETE FROM user_group_roles WHERE user_id = ? AND group_id = ?`

	// getGroupsForUser retrieves all groups a user is a member of, along with their role in each.
	getGroupsForUser = `
		SELECT g.id, g.name, g.description, ugr.role
		FROM groups g
		JOIN user_group_roles ugr ON g.id = ugr.group_id
		WHERE ugr.user_id = ?`

	// getUsersInGroup retrieves all users who are members of a specific group, along with their roles.
	getUsersInGroup = `
		SELECT u.id, u.username, ugr.role
		FROM users u
		JOIN user_group_roles ugr ON u.id = ugr.user_id
		WHERE ugr.group_id = ?`

	// getUserRoleInGroup gets a user's specific role within a single group.
	getUserRoleInGroup = `SELECT role FROM user_group_roles WHERE user_id = ? AND group_id = ?`

	// --- Relationship Queries (video_groups table) ---

	// shareVideoWithGroup makes a video accessible to members of a group.
	// 'INSERT OR IGNORE' prevents errors if the share link already exists.
	shareVideoWithGroup = `INSERT OR IGNORE INTO video_groups (video_id, group_id) VALUES (?, ?)`

	// unshareVideoFromGroup revokes a group's access to a video.
	unshareVideoFromGroup = `DELETE FROM video_groups WHERE video_id = ? AND group_id = ?`

	// getGroupsForVideo retrieves all groups a specific video is shared with.
	getGroupsForVideo = `
		SELECT g.id, g.name
		FROM groups g
		JOIN video_groups vg ON g.id = vg.group_id
		WHERE vg.video_id = ?`

	// --- Complex Access Control Query ---

	// getVisibleVideosForUser is the main query for a logged-in user. It retrieves all videos they can see:
	// 1. Videos they personally uploaded.
	// 2. Videos shared with any group of which they are a member.
	// 'DISTINCT' is used to ensure a video isn't listed twice if a user uploaded it AND it was shared with their group.
	getVisibleVideosForUser = `
		SELECT DISTINCT v.id, v.filepath, v.name, v.uploader_user_id, v.created_at
		FROM videos v
		LEFT JOIN video_groups vg ON v.id = vg.video_id
		LEFT JOIN user_group_roles ugr ON vg.group_id = ugr.group_id
		WHERE v.uploader_user_id = ? OR ugr.user_id = ?
		ORDER BY v.created_at DESC`
)
