package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

// Constants for format checking
const (
	TargetFormat  = "mov,mp4,m4a,3gp,3g2,mj2" // The format names returned by ffprobe for standard web containers
	FailingFormat = "mpegts"                  // The container format that causes web playback issues
)

// CheckVideoFormat runs ffprobe on the given file path to determine its container format.
// It returns true if the format is the problematic MPEG-TS, and the detected format name.
func CheckVideoFormat(filePath string) (bool, string, error) {
	// ffprobe command to extract the format name
	// -v error: suppress verbose output, only show errors
	// -show_entries format=format_name: only show the format name
	// -of default=noprint_wrappers=1:nokey=1: format the output to be just the raw value
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=format_name",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	)

	// Execute the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If ffprobe fails (e.g., file not found, bad file), return the error
		return false, "", fmt.Errorf("ffprobe execution failed for %s: %w\nOutput: %s", filePath, err, string(output))
	}

	// The output is the format name, possibly with surrounding whitespace
	formatName := strings.TrimSpace(string(output))

	// Check if the detected format is the known failing format (mpegts)
	needsFix := strings.EqualFold(formatName, FailingFormat)

	return needsFix, formatName, nil
}

// FixMpegTsToMp4 re-muxes the input file from MPEG-TS to MP4 container using ffmpeg.
// It uses -c copy to avoid re-encoding and -movflags faststart for web optimization.
func FixMpegTsToMp4(inputPath, outputPath string) error {
	fmt.Printf("--- Fixing %s -> %s (Re-muxing)\n", inputPath, outputPath)

	// ffmpeg command for lossless container conversion
	// -i: input file
	// -c copy: copy all streams (video/audio) without re-encoding (fast and lossless)
	// -movflags faststart: optimize the MP4 for web streaming
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-c", "copy",
		"-movflags", "faststart",
		outputPath,
	)

	// Run the command and pipe the output to a buffer to capture errors
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg re-muxing failed for %s: %w\nOutput: %s", inputPath, err, string(output))
	}

	fmt.Printf("--- Successfully re-muxed to %s\n", outputPath)
	return nil
}
 
func VerifyCodec(pathToCheck string) {
	needsFix, format, err := CheckVideoFormat(pathToCheck)
	if err != nil {
		log.Printf("Error during format check: %v\n", err)
		if strings.Contains(err.Error(), "executable file not found") {
			log.Println("'ffprobe' was not found. Please ensure ffmpeg/ffprobe is installed and in your system PATH.")
		} else if strings.Contains(err.Error(), "No such file or directory") || strings.Contains(err.Error(), "Invalid argument") {
			log.Println("The file path used in the demo is hypothetical. Please change 'pathToCheck' in main() to a real video file path.")
		}
		os.Exit(1)
	}
	if needsFix {
		base := strings.TrimSuffix(pathToCheck, filepath.Ext(pathToCheck))
		fixedPath := fmt.Sprintf("%s_fixed_%d.mp4", base, time.Now().Unix())
		fixErr := FixMpegTsToMp4(pathToCheck, fixedPath)
		if fixErr != nil {
			fmt.Printf("FATAL Error during fix process: %v\n", fixErr)
		} else {
			fmt.Printf("SUCCESS: Video successfully converted to web-compatible MP4 container at: %s\n", fixedPath)
		}

	} else {
		if strings.Contains(format, "mp4") || strings.Contains(format, "mov") {
			return
		} else {
			log.Printf("Unknown video format. Detected '%s'. Video might not be playable in the web browser\n", format)
		}
	}
}
