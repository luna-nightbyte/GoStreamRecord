package db

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"remoteCtrl/internal/utils"
	"strconv"
	"time"
)

// AddVideo inserts a new video record.
func (v *Video) Add(ctx context.Context, videoFilepath string, user_id int) error {
	if videoFilepath == "" {
		return errors.New("video filepath and downloader username cannot be empty")
	}
	now := time.Now().Format(time.RFC3339)
	videoName := filepath.Base(videoFilepath)

	sha256, _ := utils.FileSHA256(videoFilepath)
	_, err := DataBase.SQL.ExecContext(ctx, createVideo, videoFilepath, videoName, sha256, user_id, now)

	if err != nil {
		return err
	}

	return nil
}

// AddVideo inserts a new video record.
func (v *Video) Share(videoID, groupID int) error {

	_, err := DataBase.SQL.ExecContext(DataBase.ctx, shareVideoWithGroup, videoID, groupID)
	return err
}
func (v *Video) NameToID(name string) int {
	r, _ := v.ListAll(DataBase.ctx)
	return r[name].ID

}

// ListAllVideos retrieves all videos from the database.
func (v *Video) ListAll(ctx context.Context) (map[string]Video, error) {
	usr_id := strconv.Itoa(DataBase.Users.NameToID(InternalUser))
	rows, err := DataBase.SQL.QueryContext(ctx, getVisibleVideosForUser, usr_id, usr_id)
	if err != nil {
		return nil, fmt.Errorf("failed to query videos: %w", err)
	}
	defer rows.Close()

	videoMap := make(map[string]Video)
	for rows.Next() {
		var v Video
		var updatedAt string
		if err := rows.Scan(&v.ID, &v.Filepath, &v.Name, &v.Sha256, &v.UploaderUserID, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan video row: %w", err)
		}
		if v.CreatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
			return nil, fmt.Errorf("failed to parse timestamp for video %s: %w", v.Name, err)
		}
		videoMap[v.Filepath] = v
	}

	return videoMap, rows.Err()
}

// In db/video.go

func (v *Video) ListAvailable(ctx context.Context, userID int) ([]Video, error) {

	rows, err := DataBase.SQL.QueryContext(ctx, getVisibleVideosForUser, userID, userID)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to query visible videos: %w", err)
	}
	defer rows.Close()

	// Using a map to handle potential duplicates from the UNION
	videoMap := make(map[int]Video)
	for rows.Next() {
		var v Video
		var createdAt string
		if err := rows.Scan(&v.ID, &v.Filepath, &v.Name, &v.Sha256, &v.UploaderUserID, &createdAt); err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("failed to scan video row: %w", err)
		}
		v.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		videoMap[v.ID] = v
	}

	// Convert map to slice
	videos := make([]Video, 0, len(videoMap))
	for _, video := range videoMap {
		videos = append(videos, video)
	}

	return videos, rows.Err()
}

// UserHasAccessToVideo checks if a user has access to a specific video.
// This is a simplified check based on who downloaded it.
// A more robust implementation would check against user groups.
func (v *Video) CheckUserAccess(ctx context.Context, username string, videoName string) (bool, error) {
	query := "SELECT COUNT(*) FROM videos WHERE downloaded_by = ? AND name = ?"
	var count int
	err := DataBase.SQL.QueryRowContext(ctx, query, username, videoName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to query video access: %w", err)
	}
	return count > 0, nil
}
