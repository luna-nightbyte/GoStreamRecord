package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"remoteCtrl/internal/system/prettyprint"
	"remoteCtrl/internal/utils"
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

var exampleAdmin, defaultPass string = "admin", "password"
var exampleViewer string = "viewer"
var exampleMod string = "mod"
var randPass, _ = hashPassword(utils.RandString(15))

// createSchema executes the necessary SQL statements to build the database tables.
func createSchema(ctx context.Context, db *sql.DB) error {
	schemaSQL := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s",
		q_create_users,
		q_create_goups,
		q_create_tabs,
		q_create_user_group_roles,
		q_create_videos,
		q_create_video_groups,
		q_create_tab_groups,
		q_create_streamers,
		q_create_streamer_groups,
		q_create_apis,
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

		if err := DataBase.Groups.New(GroupAdmins, "admins with full control"); err != nil {
			log.Fatalf("Fatal: Could not create default group: %v", err)
		}

		if err := DataBase.Groups.New(GroupViewerOnly, "only viewing content"); err != nil {
			log.Fatalf("Fatal: Could not create default group: %v", err)
		}
		if err := DataBase.Groups.New(GroupDownloadAndView, "downloading and viewing content"); err != nil {
			log.Fatalf("Fatal: Could not create default group: %v", err)
		}
		// Users
		if err := DataBase.Users.New(exampleAdmin, defaultPass); err != nil {
			log.Fatalf("Fatal: Could not add default admin user: %v", err)
		}
		if err := DataBase.Users.New(exampleViewer, exampleViewer); err != nil {
			log.Fatalf("Fatal: Could not add default admin user: %v", err)
		}
		if err := DataBase.Users.New(exampleMod, exampleMod); err != nil {
			log.Fatalf("Fatal: Could not add default admin user: %v", err)
		}
		// Add to group

		// -- main user
		viewer_id := DataBase.Users.NameToID(exampleViewer)
		viewer_group_id := DataBase.Groups.NameToID(GroupViewerOnly)
		DataBase.Groups.AddUser(viewer_id, viewer_group_id, RoleUsers)

		// -- example moderator user
		mod_id := DataBase.Users.NameToID(exampleMod)
		mod_group_id := DataBase.Groups.NameToID(GroupDownloadAndView)
		DataBase.Groups.AddUser(mod_id, mod_group_id, RoleUsers)

		// -- example viewer user
		admin_id := DataBase.Users.NameToID(exampleAdmin)
		admin_group_id := DataBase.Groups.NameToID(GroupAdmins)
		DataBase.Groups.AddUser(admin_id, admin_group_id, RoleAdmin)

		// -- internal server user
		DataBase.Users.New(InternalUser, string(randPass))
		internalID := DataBase.Users.NameToID(InternalUser)
		DataBase.Groups.AddUser(internalID, admin_group_id, RoleAdmin)

		prettyprint.P.Yellow.Println("New database created:")
		prettyprint.P.BoldWhite.Println("	User:	| Password:")
		prettyprint.P.BoldGrey.Print("	", exampleAdmin)
		prettyprint.P.BoldWhite.Print("	| ")
		prettyprint.P.FaintWhite.Println(defaultPass)
		prettyprint.P.BoldGrey.Print("	", exampleMod)
		prettyprint.P.BoldWhite.Print("	| ")
		prettyprint.P.FaintWhite.Println(exampleMod)
		prettyprint.P.BoldGrey.Print("	", exampleViewer)
		prettyprint.P.BoldWhite.Print("	| ")
		prettyprint.P.FaintWhite.Println(exampleViewer)
		// TABS --------------------------------------------

		err := DataBase.Tabs.New(TabDownload, "Download videos directly from websites")
		if err != nil {
			log.Fatalf("Fatal: Could not create tab: %v", err)
		}
		err = DataBase.Tabs.New(TabGallery, "View downloaded videos and recordings")
		if err != nil {
			log.Fatalf("Fatal: Could not create tab: %v", err)
		}
		err = DataBase.Tabs.New(TabLiveStream, "Watch models live")
		if err != nil {
			log.Fatalf("Fatal: Could not create tab: %v", err)
		}
		err = DataBase.Tabs.New(TabRecorder, "Record videos from livestreams")
		if err != nil {
			log.Fatalf("Fatal: Could not create tab: %v", err)
		}
		tabs, err := DataBase.Tabs.List()
		if err != nil {
			log.Fatalf("Fatal: Could not create tab: %v", err)
		}

		// Share all with admins and mods
		for _, tab := range tabs {
			err = DataBase.Tabs.ShareTab(tab.ID, admin_group_id)
			if err != nil {
				fmt.Println("Fatal: Could not create tab: %v", err)
			}
			err = DataBase.Tabs.ShareTab(tab.ID, mod_group_id)
			if err != nil {
				fmt.Println("Fatal: Could not create tab: %v", err)
			}
		}
		err = DataBase.Tabs.ShareTab(tabs[TabGallery].ID, viewer_group_id)
		if err != nil {
			fmt.Println("Fatal: Could not create tab: %v", err)
		}
		err = DataBase.Tabs.ShareTab(tabs[TabLiveStream].ID, viewer_group_id)
		if err != nil {
			fmt.Println("Fatal: Could not create tab: %v", err)
		}
		DataBase.NewStreamer("test-streamer", "chaturbate", DataBase.Users.NameToID(exampleAdmin))

		cfg, _ := DataBase.Config()
		cfg.Port = 8050
		cfg.OutputFolder = "videos"
		cfg.EnableTelegram = false
		cfg.EnableGDrive = false
		cfg.EnableRateLimit = false
		DataBase.SaveConfig(cfg)
	}
}

func (s *DB) NewStreamer(streamer_name, provider string, user_id int) {
	groups, _, _ := DataBase.Groups.ListGroupsByUserID(user_id)
	DataBase.Streamers.New(streamer_name, provider)
	streamers, _ := DataBase.Streamers.List()
	for _, group := range groups {
		DataBase.Streamers.Share(streamers[streamer_name].ID, group.ID)

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
