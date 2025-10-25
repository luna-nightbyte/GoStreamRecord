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
func (db *DB) AddVideo(videoFilepath string, user_id int) error {
	if videoFilepath == "" {
		return errors.New("video filepath and downloader username cannot be empty")
	}
	now := time.Now().Format(time.RFC3339)
	videoName := filepath.Base(videoFilepath)

	sha256, _ := utils.FileSHA256(videoFilepath)
	_, err := db.SQL.ExecContext(db.ctx, createVideo, videoFilepath, videoName, sha256, user_id, now)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) ShareVideo(videoID, groupID int) error {
	_, err := db.SQL.ExecContext(db.ctx, shareVideoWithGroup, videoID, groupID)
	return err
}
func (db *DB) VideoNameToID(name string) int {
	r, _ := db.ListVideos()
	return r[name].ID

}

// ListAllVideos retrieves all videos from the database.
func (db *DB) ListVideos() (map[string]Video, error) {
	usr_id := strconv.Itoa(db.UserNameToID(InternalUser))
	rows, err := db.SQL.QueryContext(db.ctx, getVisibleVideosForUser, usr_id, usr_id)
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

func (db *DB) ListAvailableVideosForUser(userID int) ([]Video, error) {

	rows, err := db.SQL.QueryContext(db.ctx, getVisibleVideosForUser, userID, userID)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to query visible videos: %w", err)
	}
	defer rows.Close()

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
func (db *DB) CheckUserVideoAccess(ctx context.Context, username string, videoName string) (bool, error) {
	//query := "SELECT COUNT(*) FROM videos WHERE downloaded_by = ? AND name = ?"
	var count int
	userID := db.UserNameToID(username)
	rows, err := db.SQL.QueryContext(db.ctx, getVisibleVideosForUser, userID, userID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	db.ListVideos()
	for rows.Next() {
		var v Video
		var createdAt string
		if err := rows.Scan(&v.ID, &v.Filepath, &v.Name, &v.Sha256, &v.UploaderUserID, &createdAt); err != nil {
			fmt.Println(err)
			return false, fmt.Errorf("failed to scan video row: %w", err)
		}

		if v.Name == videoName {
			if v.UploaderUserID == userID {

			}
		}
	}
	return count > 0, nil
}
