package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GetVideoDuration(videoPath string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries",
		"format=duration", "-of", "default=noprint_wrappers=1:nokey=1", videoPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, err
	}
	durationStr := strings.TrimSpace(out.String())
	return strconv.ParseFloat(durationStr, 64)
}
func ExtractFrame(videoPath string) (string, error) {
	// Create a temporary file with .jpg extension
	tempFile, err := os.CreateTemp("", "frame-*.jpg")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	outputPath := tempFile.Name()
	duration, err := GetVideoDuration(videoPath)
	if err != nil {
		return "", err
	}
	middle := fmt.Sprintf("%.2f", duration/2)

	cmd := exec.Command("ffmpeg", "-y", "-ss", middle, "-i", videoPath, "-frames:v", "1", outputPath)

	//cmd := exec.Command("ffmpeg", "-y", "-ss", "00:00:01", "-i", videoPath, "-frames:v", "1", outputPath)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout  
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg command failed: %w", err)
	}

	return outputPath, nil
}
