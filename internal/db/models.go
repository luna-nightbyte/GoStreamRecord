package db

import "time"

// Models & Table Structures

type Config struct {
	Port int `json:"port"`

	// Web request rate limiting
	EnableRateLimit bool   `json:"enable_rate_limit"`
	RateLimit       int    `json:"rate_limit"`
	OutputFolder    string `json:"output_folder"`

	// Online staus ticker
	LoopInterval int `json:"online_check_min_ticker"`

	EnableGDrive   bool   `json:"enable_google_drive"`
	GDriveFilepath string `json:"google_drive_path"`

	TelegramChatID string `json:"chatID"`
	TelegramToken  string `json:"token"`
	EnableTelegram bool   `json:"enable_telegram"`
}

type telegram struct {
}

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

type Tab struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type Streamer struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Provider       string `json:"provider"`
	UploaderUserID int    `json:"uploader_user_id"`
}

type Api struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	OwnerID int    `json:"owner_id"`
	Key     string `json:"key"`
	Expires string `json:"expires"`
	Created string `json:"created"`
}

type user_group_relations struct {
	UserID  int    `json:"user_id"`
	GroupID string `json:"group_id"`
	Role    string `json:"role"`
}

type tab_group_relations struct {
	TabID   int    `json:"tab_id"`
	GroupID string `json:"group_id"`
}

type streamer_group_relations struct {
	StreamerID int    `json:"streamer_id"`
	GroupID    string `json:"group_id"`
}

type Video struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Sha256         string    `json:"sha256"`
	Filepath       string    `json:"filepath"`
	UploaderUserID int       `json:"uploader_user_id"`
	CreatedAt      time.Time `json:"created_at"`
}

// Roles define broad levels of access.
const (
	GroupAdmins          string = "admins"
	GroupViewerOnly      string = "viewer"
	GroupDownloadAndView string = "mod"
	RoleAdmin            string = "admin"
	RoleUsers            string = "user"
)

const (
	TabGallery    string = "gallery_tab"
	TabDownload   string = "download_tab"
	TabLiveStream string = "live_tab"
	TabRecorder   string = "recorder_tab"
	TabSettings   string = "settings_tab"
	TabAbout      string = "about_tab"
)

// initial queries
const (
	q_create_users string = `CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash BLOB NOT NULL,
    created_at TEXT NOT NULL
);`
	q_create_goups string = `CREATE TABLE groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TEXT NOT NULL
);`
	q_create_tabs string = `CREATE TABLE tabs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT
);`
	q_create_streamers string = `CREATE TABLE streamers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    uploader_user_id INTEGER NOT NULL,
    provider TEXT,
    FOREIGN KEY (uploader_user_id) REFERENCES users (id)
);`
	q_create_streamer_groups string = `CREATE TABLE streamer_group_relations (
	streamer_id INTEGER NOT NULL,
	group_id INTEGER NOT NULL,
	PRIMARY KEY (streamer_id, group_id),
	FOREIGN KEY (streamer_id) REFERENCES streamers (id) ON DELETE CASCADE,
	FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
);`
	q_create_apis string = `CREATE TABLE apis (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	owner_id TEXT NOT NULL,
	key TEXT NOT NULL,
	expires TEXT, -- RFC3339
	created TEXT NOT NULL
);`
	q_create_videos string = `CREATE TABLE videos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    filepath TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
	sha256 TEXT NOT NULL,
    uploader_user_id INTEGER NOT NULL,
    created_at TEXT NOT NULL,
    FOREIGN KEY (uploader_user_id) REFERENCES users (id)
);`
	q_create_user_group_roles string = `CREATE TABLE user_group_roles (
    user_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    role TEXT NOT NULL,  
    PRIMARY KEY (user_id, group_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (group_id) REFERENCES groups (id)
);`
	q_create_video_groups string = `CREATE TABLE video_groups (
    video_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    PRIMARY KEY (video_id, group_id),
    FOREIGN KEY (video_id) REFERENCES videos (id),
    FOREIGN KEY (group_id) REFERENCES groups (id)
);`
	q_create_tab_groups string = `CREATE TABLE tab_group_relations (
    tab_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    PRIMARY KEY (tab_id, group_id),
    FOREIGN KEY (tab_id) REFERENCES tabs (id),
    FOREIGN KEY (group_id) REFERENCES groups (id)
);`
	q_create_config string = `CREATE TABLE IF NOT EXISTS config (
    id INTEGER PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    port INTEGER,
    enable_rate_limit INTEGER,
    rate_limit INTEGER,
    output_folder TEXT,
    online_check_min_ticker INTEGER,
    enable_google_drive INTEGER,
    google_drive_path TEXT,
    chat_id TEXT,
    token TEXT,
    enable_telegram INTEGER
);`
)

// Reusable queries
const (
	select_shared_videos string = ` 
	SELECT v.id, v.name, v.filepath
	FROM videos v
	JOIN video_groups vg ON v.id = vg.video_id
	JOIN user_group_roles ugr ON vg.group_id = ugr.group_id
	WHERE ugr.user_id = ? -- The logged-in user's ID
 
	SELECT v.id, v.name, v.filepath
	FROM videos v
	WHERE v.uploader_user_id = ?; -- The logged-in user's ID`

	admin_del_user string = `DELETE FROM users WHERE id = ?;`
	get_users      string = `SELECT user_id, group_id, role FROM user_group_roles WHERE user_id = ?`

	mod_share_video string = `INSERT INTO video_groups (video_id, group_id) VALUES (?, ?);` // video_id, group_id
)

//  GROUPS ----------------------------------------------------------------------------------------

// This file centralizes all SQL queries for maintainability.
// Using constants for queries makes the Go code cleaner and prevents typos.

const (
	saveConfig = `INSERT INTO config (
    id, port, enable_rate_limit, rate_limit, output_folder, 
    online_check_min_ticker, enable_google_drive, google_drive_path, 
    chat_id, token, enable_telegram
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    port = excluded.port,
    enable_rate_limit = excluded.enable_rate_limit,
    rate_limit = excluded.rate_limit,
    output_folder = excluded.output_folder,
    online_check_min_ticker = excluded.online_check_min_ticker,
    enable_google_drive = excluded.enable_google_drive,
    google_drive_path = excluded.google_drive_path,
    chat_id = excluded.chat_id,
    token = excluded.token,
    enable_telegram = excluded.enable_telegram;`

	getConfig = `SELECT 
    port, enable_rate_limit, rate_limit, output_folder, 
    online_check_min_ticker, enable_google_drive, google_drive_path, 
    chat_id, token, enable_telegram 
FROM config WHERE id = 1;`
	// --- API Queries (apis table) ---

	// createUser inserts a new user record.
	createApi = `INSERT INTO apis (name, owner_id, key, expires, created) VALUES (?, ?, ?, ?, ?)`
	// getApiByID retrieves a single user with their password hash for authentication.
	listApis = `SELECT id, name, owner_id, key, expires, created FROM apis ORDER BY id`

	getUserApis = `
		SELECT DISTINCT a.id, a.name, a.key, a.expires, a.created
		FROM apis a 
		WHERE a.owner_id = ? 
		ORDER BY a.id DESC`
	// deleteApi removes api from user. Cascading deletes will handle their relationships.
	deleteApi = `DELETE FROM apis WHERE owner = ? AND id = ?`
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

	getAllGroups = `SELECT id, name, description, created_at FROM groups ORDER BY id DESC`

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
	createVideo = `INSERT INTO videos (filepath, name, sha256, uploader_user_id, created_at) VALUES (?, ?, ?, ?, ?)`

	// getVideoByID retrieves a single video by its ID.
	getVideoByID = `SELECT id, filepath, name, uploader_user_id, created_at FROM videos WHERE id = ?`

	// listAllVideos retrieves all videos, useful for an admin view.
	listAllVideos = `SELECT id, filepath, name, uploader_user_id, created_at FROM videos ORDER BY name`

	// deleteVideo removes a video. Cascading deletes will handle its shares.
	deleteVideo = `DELETE FROM videos WHERE id = ?`

	// --- Relationship Queries (user_group_roles table) ---

	// addUserToGroup assigns a user to a group with a specific role.
	// 'INSERT OR REPLACE' is used to easily change a user's role.
	addUserToGroup        = `INSERT OR REPLACE INTO user_group_roles (user_id, group_id, role) VALUES (?, ?, ?)`
	getUserGroupRelations = `SELECT user_id, group_id, role FROM user_group_roles WHERE user_id = ?`

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

	// --- Tabs Queries ---

	// createTab inserts a new user record.
	createTab = `INSERT INTO tabs (name, description) VALUES (?, ?)`

	// addTabToGroup = `INSERT OR REPLACE INTO user_group_roles (user_id, tab_id) VALUES (?, ?)`
	// listUsers retrieves all users without their password hashes for general listings.
	listTabs = `SELECT id, name, description FROM tabs ORDER BY id`

	// getVisibleVideosForUser is the main query for a logged-in user. It retrieves all videos they can see:
	// 1. Videos they personally uploaded.
	// 2. Videos shared with any group of which they are a member.
	// 'DISTINCT' is used to ensure a video isn't listed twice if a user uploaded it AND it was shared with their group.
	getVisibleVideosForUser = `
		SELECT DISTINCT v.id, v.filepath, v.name, v.sha256, v.uploader_user_id, v.created_at
		FROM videos v
		LEFT JOIN video_groups vg ON v.id = vg.video_id
		LEFT JOIN user_group_roles ugr ON vg.group_id = ugr.group_id
		WHERE v.uploader_user_id = ? OR ugr.user_id = ?
		ORDER BY v.id DESC`

	// 'INSERT OR IGNORE' prevents errors if the share link already exists.
	shareTabWithGroup = `INSERT OR IGNORE INTO tab_group_relations (tab_id, group_id) VALUES (?, ?)`

	// unshareVideoFromGroup revokes a group's access to a video.
	unshareTabFromGroup = `DELETE FROM tab_group_relations WHERE tab_id = ? AND group_id = ?`

	getVisibleTabsForUser = `
        SELECT DISTINCT t.id, t.name, t.description
        FROM tabs t
        JOIN tab_group_relations tg ON t.id = tg.tab_id
        JOIN user_group_roles ugr ON tg.group_id = ugr.group_id
        WHERE ugr.user_id = ?
        ORDER BY t.id`

	// --- Streamers Queries ---

	// createTab inserts a new user record.
	createStreamer = `INSERT INTO streamers (name, provider, uploader_user_id) VALUES (?, ?, ?)`

	// addTabToGroup = `INSERT OR REPLACE INTO user_group_roles (user_id, tab_id) VALUES (?, ?)`
	// listUsers retrieves all users without their password hashes for general listings.
	listStreamer = `SELECT id, name, provider, uploader_user_id FROM streamers ORDER BY id`

	// 'INSERT OR IGNORE' prevents errors if the share link already exists.
	shareStreamerbWithGroup = `INSERT OR IGNORE INTO streamer_group_relations (streamer_id, group_id) VALUES (?, ?)`

	// unshareVideoFromGroup revokes a group's access to a video.
	unshareStreamerFromGroup       = `DELETE FROM streamer_group_relations WHERE streamer_id = ? AND group_id = ?`
	removeUploaderUserFromStreamer = `DELETE FROM streamers WHERE id = ? AND uploader_user_id = ?`

	getVisibleStreamerForUser = `
        SELECT DISTINCT s.id, s.name, s.provider, s.uploader_user_id
		FROM streamers s
		JOIN streamer_group_relations sg ON s.id = sg.streamer_id
		JOIN user_group_roles ugr ON sg.group_id = ugr.group_id
		WHERE ugr.user_id = ?
		ORDER BY s.id;`

	getVisibleStreamerForGroup = `
        SELECT DISTINCT s.id, s.name, s.provider
		FROM streamers s
		JOIN streamer_group_relations sg ON s.id = sg.streamer_id 
		WHERE sg.group_id = ?
		ORDER BY s.id;`
)
