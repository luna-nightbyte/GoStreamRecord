# db 

## Global functions
```go

// ---- CONFIG 

// SaveConfig saves or updates the single application configuration row (id=1).
// This function uses an "UPSERT" pattern.
func (db *DB) SaveConfig(cfg Config) error 
// GetConfig retrieves the single application configuration row (id=1).
func (db *DB) Config() (Config) 


// -- Users 
func (db *User) New(username, raw_password string) error 
func (db *User) Authenticate(username, password string) (bool, error)
func (db *User) HttpRequestID(r *http.Request) int
func (db *User) GetUserByName(username string) (*User, error)
func (db *User) GetUserByID(id int) (*User, error)
func (db *User) IsAdmin(username string) (bool, error) 
func (db *User) Update(userID int, newUsername string, newPassword string) error 
func (db *User) Delete(userID int) error 
func (db *User) List() (map[string]User, error)

// -- Videos

func (db *DB) AddVideo(ctx context.Context, videoFilepath string, downloadedBy string) error 
func (db *DB) ShareVideo(videoID, groupID int) error 
func (u *User) NameToID(name string) int
func (g *Group) NameToID(name string) int
func (db *DB) VideoNameToID(name string) int
func (db *DB) ListAllVideos(ctx context.Context) (map[string]Video, error)
func (db *DB) ListVisibleVideosForUser(ctx context.Context, userID int) ([]Video, error)
func (db *DB) UserHasAccessToVideo(ctx context.Context, username string, videoName string) (bool, error)

// -- Groups 
func (db *Group) New(groupName string, description string) error 
func (db *Group) AddUser(userID, groupID int, role string) error
func (db *Group) ListGroupsByUserID(user_id int) (map[string]Group, string, error)
func (db *User) GetUserGroupRelations(user_id int) (user_group_relations, error)
func (db *User) GetGroupByName(username string) (*User, error)
func (db *Group) List() (map[string]Group, error) 

// -- Tabs
func (db *Tab) New(tabName, description string) error 
func (db *Tab) GetAvailableTabsForUser(userID int) (map[string]Tab, error) 
func (db *Tab) DeleteForGroup(groupID, tabID int) error
func (db *Tab) ShareTab(tabID, groupID int) error
func (db *Tab) List() (map[string]Tab, error) 


// -- Streamers
func (db *Streamer) New(streamerName, provider string) error
func (db *Streamer) Share(streamerID, groupID int) error 
func (db *Streamer) List() (map[string]Streamer, error) 
func (db *Streamer) GetAvailableForUser(userID int) (map[string]Streamer, error)
func (db *Streamer) GetAvailableForGroup(groupID int) (map[string]Streamer, error)
func (db *Streamer) DeleteForUser(user_id, streamer_id int) (*Streamer, error) 
func (db *Streamer) DeleteForGroup(groupID, streamerID int) error 
func (db *Streamer) Share(streamerID, groupID int) error


// -- API 
func (db *Api) New(apiName, username string) error
func (db *Api) ListUserApis(user_id int) (map[string]Api, error)
func (db *Api) List(owner_id int) (map[string]Api, error) 
func (db *Api) DeleteForUser(user_id, api_id int) (*Api, error)
```

## Example
```go
package main

import (
  "context"
  "log"
  "ProjectName/internal/db"
)

func main() {
  ctx := context.Background()

  // Pass "" to use default ./db/database.sqlite or set DB_PATH env var
  db.Init(ctx, "")

  // Create a user and put them in the "viewers" group
  if err := db.DataBase.Users.New("alice", "viewers"); err != nil { log.Fatal(err) }
  aliceID := db.DataBase.Users.NameToID("alice")
  viewID := db.DataBase.Groups.NameToID(db.GroupViewerOnly)
  if err := db.DataBase.Groups.AddUser(aliceID, modID, db.RoleUsers); err != nil { log.Fatal(err) }

  // Create a tab and share with the mod group
  if err := db.DataBase.Tabs.New("gallery-tab", "View video files"); err != nil { log.Fatal(err) }
  tabs, _ := db.DataBase.Tabs.List()
  if err := db.DataBase.Tabs.ShareTab(tabs["gallery-tab"].ID, viewID); err != nil { log.Fatal(err) }


  if err := db.DataBase.Users.New("peter", "admins"); err != nil { log.Fatal(err) }
  peterID := db.DataBase.Users.NameToID("peter")
  adminID := db.DataBase.Groups.NameToID(db.GroupAdmin)
  if err := db.DataBase.Groups.AddUser(peterID, adminID, db.RoleAdmins); err != nil { log.Fatal(err) }

  // Create a tab and share with the mod group
  if err := db.DataBase.Tabs.New("gallery-tab", "View video files"); err != nil { log.Fatal(err) }
  tabs, _ := db.DataBase.Tabs.List()
  if err := db.DataBase.Tabs.ShareTab(tabs["gallery-tab"].ID, viewID); err != nil { log.Fatal(err) }

  // List tabs visible to Alice (only gallery tab in this instance)
  visible, _ := db.DataBase.Tabs.GetAvailableTabsForUser(aliceID)
  log.Printf("Alice can see %d tabs", len(visible))

  // List tabs visible to Peter (all 4 tabs in this instance)
  visible, _ := db.DataBase.Tabs.GetAvailableTabsForUser(aliceID)
  log.Printf("Peter can see %d tabs", len(visible))
}
```

---

## Initialization & defaults

```go
db.Init(ctx, path)
```

- Uses `DB_PATH` env var or `./db/database.sqlite` if empty.
- Creates all tables on first run, then seeds:
  - Groups: `admins`, `viewer`, `mod`.
  - Users: `admin`, `mod`, `viewer` (all with password `password`).
  - An internal server user `_internal` with a random password.
  - Tabs: `download_tab`, `gallery_tab`, `live_tab`, `recorder_tab`.
  - Shares: all tabs to `admins` and `mod`; `gallery_tab` and `live_tab` to `viewer`.
  - A sample streamer `test-streamer` shared to the creator's groups.
  - A default `Config` row (port 8050, output folder `videos`, Telegram/GDrive disabled, rate limiting disabled).

You’ll interact with the global handle:

```go
var DataBase *db.DB
```

> The package sets `DataBase = &DB{ctx, SQL, ...}` inside `Init`.

---

## Models at a glance

- **User**: `ID, Username, PasswordHash, CreatedAt`
- **Group**: `ID, Name, Description, CreatedAt`
- **Tab**: `ID, Name, Description`
- **Streamer**: `ID, Name, Provider`
- **Video**: `ID, Name, Sha256, Filepath, UploaderUserID, CreatedAt`
- **Api**: `ID, Name, Key, Expires, Created`
- **Config**: port, rate-limit flags, output folder, GDrive/Telegram settings

---

## Common tasks

### Users

```go
// Create
if err := db.DataBase.Users.New("bob", "supersecret"); err != nil { /* handle */ }

// Authenticate
ok, err := db.DataBase.Users.Authenticate("bob", "supersecret")

// List (map[username]User)
usrs, _ := db.DataBase.Users.List()

// Update username and/or password
if err := db.DataBase.Users.Update(userID, "bobby", ""); err != nil { /* handle */ }

// Delete
if err := db.DataBase.Users.Delete(userID); err != nil { /* handle */ }

// Check admin role
isAdmin, _ := db.DataBase.Users.IsAdmin("admin")

```

### Groups & roles

```go
// Create group
if err := db.DataBase.Groups.New("staff", "Company staff"); err != nil { /* handle */ }

// Add user to group (role = admin|user)
userID := db.DataBase.Users.NameToID("bob")
staffID := db.DataBase.Groups.NameToID("staff")
if err := db.DataBase.Groups.AddUser(userID, staffID, db.RoleUsers); err != nil { /* handle */ }

// List groups
all, _ := db.DataBase.Groups.List()

// Groups by user (with role)
byUser, role, _ := db.DataBase.Groups.ListGroupsByUserID(userID)
```

### Tabs (feature access)

```go
// Create a tab
_ = db.DataBase.Tabs.New("reports_tab", "View business reports")

// Share to a group
reports := db.DataBase.Tabs.List()
_ = db.DataBase.Tabs.ShareTab(reports["reports_tab"].ID, staffID)

// Visible tabs for a user
visible, _ := db.DataBase.Tabs.GetAvailableTabsForUser(userID)
```

### Streamers

```go
// Create and share
_ = db.DataBase.Streamers.New("alice_cam", "chaturbate")
streamers, _ := db.DataBase.Streamers.List()
_ = db.DataBase.Streamers.Share(streamers["alice_cam"].ID, staffID)

// Visible for user
sv, _ := db.DataBase.Streamers.GetAvailableForUser(userID)
```

### Videos

```go
// Add a record for a downloaded video
_ = db.DataBase.AddVideo(ctx, "/videos/2025-10-01/cat.mp4", "bob")

// Share with a group
vids, _ := db.DataBase.ListAllVideos(ctx)
_ = db.DataBase.ShareVideo(vids["/videos/2025-10-01/cat.mp4"].ID, staffID)

// All visible videos for a user
vv, _ := db.DataBase.ListVisibleVideosForUser(ctx, userID)
```

### API keys (per-user sharing)

```go
// Create an API key and grant to a user
_ = db.DataBase.Api.New("cli", "bob")

// List all API keys
apis, _ := db.DataBase.Api.List()

// API keys for a user
mine, _ := db.DataBase.Api.ListUserApis(userID)
```

> Note: API expiry/key rotation are simple fields; wire your own logic for generating keys and enforcing expiry.

### App config

```go
cfg, _ := db.DataBase.Config()
cfg.Port = 8080
cfg.OutputFolder = "videos"
_ = db.DataBase.SaveConfig(cfg)
```

---

## Access control model (overview)

- Users belong to **groups** with a **role** (`admin` or `user`).
- Content/features (tabs, streamers, videos) are **shared to groups**.
- "Visible for user" queries resolve visibility from (uploader OR any of the user’s groups).

---

## Environment & runtime

- **DB path**: `DB_PATH` env var or default `./db/database.sqlite`.
- **SQLite**: `SetMaxOpenConns(1)` for safety with file-based SQLite.
- **Passwords**: hashed with bcrypt.

---

## Error handling

Common errors to expect and handle in your app:

- `ErrUserNotFound` – lookups that return no rows.
- Unique constraint errors – surface as user-friendly messages like "already exists".

---

## Defaults / seed data

| Type   | Values                                                               |
| ------ | -------------------------------------------------------------------- |
| Groups | `admins` (full control), `mod` (download+view), `viewer` (view only) |
| Users  | `admin` / `mod` / `viewer` (password: `password`)                    |
| Tabs   | `download_tab`, `gallery_tab`, `live_tab`, `recorder_tab`            |
| Shares | All tabs → `admins`, `mod`; `gallery_tab`+`live_tab` → `viewer`      |
| Other  | `_internal` user (random pw), `test-streamer` example                |

> Replace or remove seed data in your own bootstrap if needed.

---

## Gotchas & tips

- When updating a user **without** a new password, pass `""` for password to keep the old hash.
- `NameToID` helpers fetch the current list each time—cache in your service layer if you’re calling them in hot paths.
- For API keys, supply your own secure key generation and expiry enforcement.

---

## Minimal table schema (for reference)

Tables: `users`, `groups`, `user_group_roles`, `tabs`, `tab_group_relations`, `streamers`, `streamer_group_relations`, `videos`, `video_groups`, `apis`, `api_user_relations`, `config`.

Migrations are all embedded in code and applied automatically on first run.

--- 
