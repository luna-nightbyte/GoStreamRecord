package db

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"
)

// AddVideo inserts a new video record.
func (db *DB) AddVideo(ctx context.Context, videoFilepath string, downloadedBy string) error {
	if downloadedBy == "" || videoFilepath == "" {
		return errors.New("video filepath and downloader username cannot be empty")
	}

	now := time.Now().Format(time.RFC3339)
	videoName := filepath.Base(videoFilepath)
	// Default to an empty list of groups, can be updated later.
	groups, err := marshalIntSlice([]int{})
	if err != nil {
		return fmt.Errorf("failed to marshal default groups for video: %w", err)
	}

	query := "INSERT INTO videos (filepath, name, downloaded_by, groups, updated_at) VALUES (?, ?, ?, ?, ?)"
	_, err = db.SQL.ExecContext(ctx, query, videoFilepath, videoName, downloadedBy, groups, now)

	if err != nil {
		return errors.New("video file already exists or a database error occurred")
	}

	return nil
}

// ListAllVideos retrieves all videos from the database.
func (db *DB) ListAllVideos(ctx context.Context) (map[string]Video, error) {
 
	rows, err := db.SQL.QueryContext(ctx, select_shared_videos)
	if err != nil {
		return nil, fmt.Errorf("failed to query videos: %w", err)
	}
	defer rows.Close()

	videoMap := make(map[string]Video)
	for rows.Next() {
		var v Video
		var groupsJSON string
		var updatedAt string
		if err := rows.Scan(&v.ID, &v.Filepath, &v.Name, &v.UploaderUserID, &groupsJSON, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan video row: %w", err)
		} 
		if v.CreatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
			return nil, fmt.Errorf("failed to parse timestamp for video %s: %w", v.Name, err)
		}
		videoMap[v.Name] = v
	}

	return videoMap, rows.Err()
}

// In db/video.go

func (db *DB) ListVisibleVideosForUser(ctx context.Context, userID int) ([]Video, error) {
	 

	rows, err := db.SQL.QueryContext(ctx, select_shared_videos, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query visible videos: %w", err)
	}
	defer rows.Close()

	// Using a map to handle potential duplicates from the UNION
	videoMap := make(map[int]Video)
	for rows.Next() {
		var v Video
		var createdAt string
		if err := rows.Scan(&v.ID, &v.Name, &v.Filepath, &v.UploaderUserID, &createdAt); err != nil {
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
func (db *DB) UserHasAccessToVideo(ctx context.Context, username string, videoName string) (bool, error) {
	query := "SELECT COUNT(*) FROM videos WHERE downloaded_by = ? AND name = ?"
	var count int
	err := db.SQL.QueryRowContext(ctx, query, username, videoName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to query video access: %w", err)
	}
	return count > 0, nil
}
