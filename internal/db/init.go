package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"remoteCtrl/internal/system/prettyprint"
	"remoteCtrl/internal/utils"
	"time"
)

// DB wraps the sql.DB connection pool.
type DB struct {
	ctx       context.Context
	SQL       *sql.DB
	Users     User
	Groups    Group
	Streamers Streamer
	Videos    Video
	Tabs      Tab
	APIs      Api
}

const default_db_path string = "./db/database.sqlite"

// Global variable to hold the database instance.
var DataBase *DB

// Internal server user
const InternalUser string = "_internal"

const exampleAdmin, defaultPass string = "admin", "password"
const exampleViewer string = "viewer"
const exampleMod string = "mod"

var randPass, _ = hashPassword(utils.RandString(15))

// createSchema executes the necessary SQL statements to build the database tables.
func createSchema(ctx context.Context, db *sql.DB) error {
	schemaSQL := fmt.Sprintf(`
	%s
	%s
	%s
	%s
	%s
	%s
	%s
	%s
	%s
	%s 
	%s`,
		q_create_users,
		q_create_groups,
		q_create_tabs,
		q_create_videos,
		q_create_apis,
		q_create_streamers,
		q_create_video_group_relations,
		q_create_tab_group_relations,
		q_create_user_group_relations,
		q_create_streamer_group_relations, 
		q_create_config,
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
		dbPath = default_db_path
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

		if err := DataBase.NewGroup(GroupAdmins, "admins with full control"); err != nil {
			log.Fatalf("Fatal: Could not create default group: %v", err)
		}

		if err := DataBase.NewGroup(GroupViewerOnly, "only viewing content"); err != nil {
			log.Fatalf("Fatal: Could not create default group: %v", err)
		}
		if err := DataBase.NewGroup(GroupDownloadAndView, "downloading and viewing content"); err != nil {
			log.Fatalf("Fatal: Could not create default group: %v", err)
		}
		// Users
		if err := DataBase.NewUser(exampleAdmin, defaultPass); err != nil {
			log.Fatalf("Fatal: Could not add default admin user: %v", err)
		}
		if err := DataBase.NewUser(exampleViewer, defaultPass); err != nil {
			log.Fatalf("Fatal: Could not add default admin user: %v", err)
		}
		if err := DataBase.NewUser(exampleMod, defaultPass); err != nil {
			log.Fatalf("Fatal: Could not add default admin user: %v", err)
		}
		if err := DataBase.NewUser(InternalUser, string(randPass)); err != nil {
			log.Fatalf("Fatal: Could not add default admin user: %v", err)
		}
		// Add to group

		// -- main user
		admin_group_id := DataBase.GroupNameToID(GroupAdmins)
		mod_group_id := DataBase.GroupNameToID(GroupDownloadAndView)
		viewer_group_id := DataBase.GroupNameToID(GroupViewerOnly)

		AddUserToGroup(exampleAdmin, GroupAdmins, RoleAdmin)
		AddUserToGroup(exampleAdmin, GroupDownloadAndView, RoleAdmin)
		AddUserToGroup(exampleAdmin, GroupViewerOnly, RoleAdmin)

		AddUserToGroup(InternalUser, GroupAdmins, RoleAdmin)
		AddUserToGroup(InternalUser, GroupDownloadAndView, RoleAdmin)
		AddUserToGroup(InternalUser, GroupViewerOnly, RoleAdmin)

		AddUserToGroup(exampleMod, GroupDownloadAndView, RoleUsers)
		AddUserToGroup(exampleMod, GroupViewerOnly, RoleAdmin)

		AddUserToGroup(exampleViewer, GroupViewerOnly, RoleUsers)

		prettyprint.P.Yellow.Println("New database created:")
		prettyprint.P.BoldWhite.Print("User:")
		prettyprint.P.BoldGrey.Println("	", exampleAdmin)
		prettyprint.P.BoldWhite.Print("User:")
		prettyprint.P.BoldGrey.Println("	", exampleMod)
		prettyprint.P.BoldWhite.Print("User:")
		prettyprint.P.BoldGrey.Println("	", exampleViewer)
		prettyprint.P.BoldWhite.Print("Password (for all): ")
		prettyprint.P.FaintWhite.Println(defaultPass)
		// TABS --------------------------------------------

		err := DataBase.NewTab(TabDownload, "Download videos directly from websites")
		if err != nil {
			log.Fatalf("Fatal: Could not create tab: %v", err)
		}
		err = DataBase.NewTab(TabGallery, "View downloaded videos and recordings")
		if err != nil {
			log.Fatalf("Fatal: Could not create tab: %v", err)
		}
		err = DataBase.NewTab(TabLiveStream, "Watch models live")
		if err != nil {
			log.Fatalf("Fatal: Could not create tab: %v", err)
		}
		err = DataBase.NewTab(TabRecorder, "Record videos from livestreams")
		if err != nil {
			log.Fatalf("Fatal: Could not create tab: %v", err)
		}

		// TODO
		err = DataBase.NewTab(TabSettings, "General settings")
		if err != nil {
			log.Fatalf("Fatal: Could not create tab: %v", err)
		}
		err = DataBase.NewTab(TabAbout, "About us")
		if err != nil {
			log.Fatalf("Fatal: Could not create tab: %v", err)
		}
		err = DataBase.NewTab(TabLogs, "System logs")
		if err != nil {
			log.Fatalf("Fatal: Could not create tab: %v", err)
		}

		tabs, err := DataBase.ListTabs()
		if err != nil {
			log.Fatalf("Fatal: Could not create tab: %v", err)
		}

		// Share all with admins and mods
		for _, tab := range tabs {
			err = DataBase.ShareTab(tab.ID, admin_group_id)
			if err != nil {
				log.Fatalf("Fatal: Could not share tab: %v", err)
			}
			err = DataBase.ShareTab(tab.ID, mod_group_id)
			if err != nil {
				log.Fatalf("Fatal: Could not share tab: %v", err)
			}
		}
		err = DataBase.ShareTab(tabs[TabGallery].ID, viewer_group_id)
		if err != nil {
			log.Fatalf("Fatal: Could not create tab: %v", err)
		}
		err = DataBase.ShareTab(tabs[TabLiveStream].ID, viewer_group_id)
		if err != nil {
			log.Fatalf("Fatal: Could not create tab: %v", err)
		}

		// DataBase.NewStreamer("test-streamer", "chaturbate", DataBase.Users.NameToID(exampleAdmin), true)

		cfg := DataBase.Config()
		cfg.Port = 8050
		cfg.OutputFolder = "videos"
		cfg.EnableTelegram = false
		cfg.EnableGDrive = false
		cfg.EnableRateLimit = false
		DataBase.SaveConfig(cfg)
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

// USER & GROUPS ------------------------------------------------------------------------------------
func (db *DB) query_row_UserSql(query string, args ...any) (User, error) {
	var usr User
	row := db.SQL.QueryRowContext(db.ctx, query, args...)
	var createdAt string
	err := row.Scan(&usr.ID, &usr.Username, &usr.PasswordHash, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return usr, ErrUserNotFound
		}
		return usr, err
	}

	if usr.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
		return usr, err
	}
	return usr, nil
}
func (db *DB) query_rows_UserGroupRelationsSql(query string, args ...any) ([]user_group_relations, error) {

	rows, err := db.SQL.QueryContext(db.ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	usrGrpRel := []user_group_relations{}
	for rows.Next() {
		var ugr user_group_relations
		if err := rows.Scan(&ugr.UserID, &ugr.GroupID, &ugr.Role); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		usrGrpRel = append(usrGrpRel, ugr)
	}

	return usrGrpRel, nil
}
