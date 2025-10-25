package db

import "time"

// Groups
const (
	GroupAdmins          string = "admins"
	GroupViewerOnly      string = "viewer"
	GroupDownloadAndView string = "mod"
)

// Roles
const (
	RoleAdmin string = "admin"
	RoleUsers string = "user"
)

// Tabs
const (
	TabGallery    string = "gallery_tab"
	TabDownload   string = "download_tab"
	TabLiveStream string = "live_tab"
	TabRecorder   string = "recorder_tab"
	TabSettings   string = "settings_tab"
	TabLogs       string = "logs_tab"
	TabAbout      string = "about_tab"
)

// Table Structures
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
	GroupID int    `json:"group_id"`
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

type video_group_relations struct {
	VideoID int    `json:"video_id"`
	GroupID string `json:"group_id"`
}

type Video struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Sha256         string    `json:"sha256"`
	Filepath       string    `json:"filepath"`
	UploaderUserID int       `json:"uploader_user_id"`
	CreatedAt      time.Time `json:"created_at"`
}

// Database structure
const (

	// User, Group, Tabs, Streamers, Apis Videos ----------------------
	q_create_users string = `CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash BLOB NOT NULL,
    created_at TEXT NOT NULL
);`
	q_create_groups string = `CREATE TABLE groups (
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

	// Relations --------------------------------------------

	q_create_user_group_relations string = `CREATE TABLE user_group_relations (
    user_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    role TEXT NOT NULL,  
    PRIMARY KEY (user_id, group_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (group_id) REFERENCES groups (id)
);`
	q_create_video_group_relations string = `CREATE TABLE video_group_relations (
    video_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    PRIMARY KEY (video_id, group_id),
    FOREIGN KEY (video_id) REFERENCES videos (id),
    FOREIGN KEY (group_id) REFERENCES groups (id)
);`
	q_create_tab_group_relations string = `CREATE TABLE tab_group_relations (
    tab_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    PRIMARY KEY (tab_id, group_id),
    FOREIGN KEY (tab_id) REFERENCES tabs (id),
    FOREIGN KEY (group_id) REFERENCES groups (id)
);`

	q_create_streamer_group_relations string = `CREATE TABLE streamer_group_relations (
	streamer_id INTEGER NOT NULL,
	group_id INTEGER NOT NULL,
	PRIMARY KEY (streamer_id, group_id),
	FOREIGN KEY (streamer_id) REFERENCES streamers (id) ON DELETE CASCADE,
	FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
);`
)

// Queries
const (

	// USERS ----------------------

	createUser        = `INSERT INTO users (username, password_hash, created_at) VALUES (?, ?, ?)`
	getUserByID       = `SELECT id, username, password_hash, created_at FROM users WHERE id = ?`
	getUserByUsername = `SELECT id, username, password_hash, created_at FROM users WHERE username = ?`
	listUsers         = `SELECT id, username, created_at FROM users ORDER BY username`
	updateUser        = `UPDATE users SET username = ?, password_hash = ? WHERE id = ?`
	deleteUser        = `DELETE FROM users WHERE id = ?`

	// APIS ----------------------

	createApi   = `INSERT INTO apis (name, owner_id, key, expires, created) VALUES (?, ?, ?, ?, ?)`
	listApis    = `SELECT id, name, owner_id, key, expires, created FROM apis ORDER BY id`
	getUserApis = `
		SELECT DISTINCT a.id, a.name, a.owner_id, a.key, a.expires, a.created
		FROM apis a 
		WHERE a.owner_id = ? 
		ORDER BY a.id DESC`
	deleteApi = `DELETE FROM apis WHERE owner = ? AND id = ?`

	// GROUPS ----------------------

	createGroup       = `INSERT INTO groups (name, description, created_at) VALUES (?, ?, ?)`
	getGroupByGroupID = `SELECT id, name, description, created_at FROM groups WHERE id = ?`
	getAllGroups      = `SELECT id, name, description, created_at FROM groups ORDER BY id DESC`
	getGroupByName    = `SELECT id, name, description, created_at FROM groups WHERE name = ?`
	listGroups        = `SELECT id, name, description, created_at FROM groups ORDER BY name`
	updateGroup       = `UPDATE groups SET name = ?, description = ? WHERE id = ?`
	deleteGroup       = `DELETE FROM groups WHERE id = ?`

	// STREAMERS ----------------------

	createStreamer = `INSERT INTO streamers (name, provider, uploader_user_id) VALUES (?, ?, ?)`
	listStreamer   = `SELECT id, name, provider, uploader_user_id FROM streamers ORDER BY id`

	// TABS ----------------------

	createTab = `INSERT INTO tabs (name, description) VALUES (?, ?)`
	listTabs  = `SELECT id, name, description FROM tabs ORDER BY id`

	// CONFG ----------------------

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

	// VIDEOS ----------------------

	createVideo   = `INSERT INTO videos (filepath, name, sha256, uploader_user_id, created_at) VALUES (?, ?, ?, ?, ?)`
	getVideoByID  = `SELECT id, filepath, name, uploader_user_id, created_at FROM videos WHERE id = ?`
	listAllVideos = `SELECT id, filepath, name, uploader_user_id, created_at FROM videos ORDER BY name`
	deleteVideo   = `DELETE FROM videos WHERE id = ?`

	// RELATIONS (aka groups) ------------------------------------------------------
	// USER & GROUPS ----------------------

	get_users_and_relations  = `SELECT user_id, group_id, role FROM user_group_relations WHERE user_id = ?`
	addUserToGroup           = `INSERT OR REPLACE INTO user_group_relations (user_id, group_id, role) VALUES (?, ?, ?)`
	getUserGroupRelations    = `SELECT user_id, group_id, role FROM user_group_relations WHERE user_id = ?`
	getAllUserGroupRelations = `SELECT user_id, group_id, role FROM user_group_relations ORDER BY user_id`
	removeUserFromGroup      = `DELETE FROM user_group_relations WHERE user_id = ? AND group_id = ?`
	getGroupsForUser         = `
		SELECT g.id, g.name, g.description, ugr.role
		FROM groups g
		JOIN user_group_relations ugr ON g.id = ugr.group_id
		WHERE ugr.user_id = ?`
	getUsersInGroup = `
		SELECT u.id, u.username, ugr.role
		FROM users u
		JOIN user_group_relations ugr ON u.id = ugr.user_id
		WHERE ugr.group_id = ?`
	getUserRoleInGroup = `SELECT role FROM user_group_relations WHERE user_id = ? AND group_id = ?`

	// VIDEOS & GROUPS/USERS ----------------------

	mod_share_video   string = `INSERT INTO video_group_relations (video_id, group_id) VALUES (?, ?);` // video_id, group_id
	get_shared_videos string = ` 
	SELECT v.id, v.name, v.filepath
	FROM videos v
	JOIN video_group_relations vg ON v.id = vg.video_id
	JOIN user_group_relations ugr ON vg.group_id = ugr.group_id
	WHERE ugr.user_id = ? -- The logged-in user's ID
 
	SELECT v.id, v.name, v.filepath
	FROM videos v
	WHERE v.uploader_user_id = ?; -- The logged-in user's ID`
	shareVideoWithGroup   = `INSERT OR IGNORE INTO video_group_relations (video_id, group_id) VALUES (?, ?)`
	unshareVideoFromGroup = `DELETE FROM video_group_relations WHERE video_id = ? AND group_id = ?`
	getGroupsForVideo     = `
		SELECT g.id, g.name
		FROM groups g
		JOIN video_group_relations vg ON g.id = vg.group_id
		WHERE vg.video_id = ?`
	getVisibleVideosForUser = `
		SELECT DISTINCT v.id, v.filepath, v.name, v.sha256, v.uploader_user_id, v.created_at
		FROM videos v
		LEFT JOIN video_group_relations vg ON v.id = vg.video_id
		LEFT JOIN user_group_relations ugr ON vg.group_id = ugr.group_id
		WHERE v.uploader_user_id = ? OR ugr.user_id = ?
		ORDER BY v.id DESC`

	// TABS & GROUPS/USERS ----------------------

	shareTabWithGroup   = `INSERT OR IGNORE INTO tab_group_relations (tab_id, group_id) VALUES (?, ?)`
	unshareTabFromGroup = `DELETE FROM tab_group_relations WHERE tab_id = ? AND group_id = ?`

	getVisibleTabsForUser = `
        SELECT DISTINCT t.id, t.name, t.description
        FROM tabs t
        JOIN tab_group_relations tg ON t.id = tg.tab_id
        JOIN user_group_relations ugr ON tg.group_id = ugr.group_id
        WHERE ugr.user_id = ?
        ORDER BY t.id`

	// STREAMERS & GROUPS/USERS ----------------------

	shareStreamerWithGroup         = `INSERT OR IGNORE INTO streamer_group_relations (streamer_id, group_id) VALUES (?, ?)`
	unshareStreamerFromGroup       = `DELETE FROM streamer_group_relations WHERE streamer_id = ? AND group_id = ?`
	removeUploaderUserFromStreamer = `DELETE FROM streamers WHERE id = ? AND uploader_user_id = ?`
	getVisibleStreamerForUser      = `
        SELECT DISTINCT s.id, s.name, s.provider, s.uploader_user_id
		FROM streamers s
		JOIN streamer_group_relations sg ON s.id = sg.streamer_id
		JOIN user_group_relations ugr ON sg.group_id = ugr.group_id
		WHERE ugr.user_id = ?
		ORDER BY s.id;`
	getVisibleStreamerForGroup = `
        SELECT DISTINCT s.id, s.name, s.provider
		FROM streamers s
		JOIN streamer_group_relations sg ON s.id = sg.streamer_id 
		WHERE sg.group_id = ?
		ORDER BY s.id;`
)
