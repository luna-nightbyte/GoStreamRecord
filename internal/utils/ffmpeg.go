package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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
func ExtractFrame(tmpDir, videoPath string) (string, error) {
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

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg command failed: %w", err)
	}

	return outputPath, nil
}

// --- Constants for Failure Detection ---
const (
	FailingContainer = "mpegts" // Problem 1: Container format
	FailingCodec     = "hevc"   // Problem 2: Video codec (H.265)
	WorkingCodec     = "h264"   // Target working video codec (AVC)
)

// VideoCheckResult stores the results of the format and codec check.
type VideoCheckResult struct {
	NeedsFix         bool
	IsContainerIssue bool
	IsCodecIssue     bool
	DetectedFormat   string
	DetectedCodec    string
}

// ffprobe JSON structure for quick parsing of format and streams
type ffprobeOutput struct {
	Format struct {
		FormatName string `json:"format_name"`
	} `json:"format"`
	Streams []struct {
		CodecName string `json:"codec_name"`
		CodecType string `json:"codec_type"`
	} `json:"streams"`
}

// CheckVideoFormat runs ffprobe and parses the output to check for both known failure modes.
func CheckVideoFormat(filePath string) (VideoCheckResult, error) {
	result := VideoCheckResult{}

	// ffprobe command to dump format and stream info as JSON
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0", // Select only the first video stream
		"-show_entries", "format=format_name:stream=codec_name,codec_type",
		"-of", "json",
		filePath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return result, fmt.Errorf("ffprobe execution failed for %s: %w\nOutput: %s", filePath, err, string(output))
	}

	var data ffprobeOutput
	if err := json.Unmarshal(output, &data); err != nil {
		return result, fmt.Errorf("failed to parse ffprobe JSON output: %w", err)
	}

	// 1. Check Container Format
	result.DetectedFormat = data.Format.FormatName
	if strings.EqualFold(result.DetectedFormat, FailingContainer) {
		result.NeedsFix = true
		result.IsContainerIssue = true
	}

	// 2. Check Video Codec
	for _, stream := range data.Streams {
		if stream.CodecType == "video" {
			result.DetectedCodec = stream.CodecName
			if strings.EqualFold(result.DetectedCodec, FailingCodec) {
				result.NeedsFix = true
				result.IsCodecIssue = true
			}
			break // Only need to check the first video stream
		}
	}

	return result, nil
}

// FixMpegTsToMp4 performs a fast, lossless container re-mux for MPEG-TS issues.
func FixMpegTsToMp4(inputPath, outputPath string) error {
	fmt.Printf("--- Fix: Lossless Re-muxing %s (MPEG-TS -> MP4)\n", inputPath)

	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-c", "copy", // Copy streams without re-encoding
		"-movflags", "faststart",
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg re-muxing failed for %s: %w\nOutput: %s", inputPath, err, string(output))
	}

	fmt.Printf("--- Successfully re-muxed to %s\n", outputPath)
	return nil
}

func FixHevcToAvc(inputPath, outputPath string) error {

	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-c:v", "libx264", // Force H.264 codec
		"-preset", "medium",
		"-crf", "23", // Constant Rate Factor: 23 is generally considered visually lossless
		"-c:a", "copy",
		"-movflags", "faststart",
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg re-encoding failed for %s: %w\nOutput: %s", inputPath, err, string(output))
	}

	fmt.Printf("--- Successfully re-encoded to web-compatible H.264 MP4 at %s\n", outputPath)
	return nil
}

func checkFFmpegDependencies() {
	for _, bin := range []string{"ffprobe", "ffmpeg"} {
		if _, err := exec.LookPath(bin); err != nil {
			fmt.Printf("\nFATAL: '%s' was not found. Please ensure ffmpeg/ffprobe is installed and in your system PATH.\n", bin)
			os.Exit(1)
		}
	}
}

type VideoIntegrity struct {
	mu    sync.Mutex
	Queue []string // queue of videos that needs to be checked.
}

var VideoVerify VideoIntegrity

func (vi *VideoIntegrity) Contains(videoPath string) bool {
	vi.mu.Lock()
	defer vi.mu.Unlock()
	for _, video := range vi.Queue {
		if video == videoPath {
			return true
		}
	}
	return false
}
func (vi *VideoIntegrity) Add(videoPath string) {
	if vi.Contains(videoPath) {
		return
	}
	vi.mu.Lock()
	defer vi.mu.Unlock()
	vi.Queue = append(vi.Queue, videoPath)
}
func (vi *VideoIntegrity) RunCodecVerification() {
	vi.mu.Lock()
	queueCopy := make([]string, len(vi.Queue))
	copy(queueCopy, vi.Queue)
	vi.Queue = nil
	vi.mu.Unlock()
	for _, video := range queueCopy {
		VerifyCodec(video)
	}
}

func VerifyCodec(pathToCheck string) {
	results, err := CheckVideoFormat(pathToCheck)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "executable file not found"):
			log.Println("ERROR: 'ffprobe' not found. Please install ffmpeg/ffprobe and ensure it's in PATH.")
		case strings.Contains(err.Error(), "No such file or directory"):
			log.Printf("ERROR: File not found: %s\n", pathToCheck)
		case strings.Contains(err.Error(), "Invalid argument"):
			log.Printf("ERROR: Invalid input file: %s\n", pathToCheck)
		default:
			log.Printf("ERROR during format check: %v\n", err)
		}
		return
	}

	if !results.NeedsFix {
		return
	}

	var fixFunc func(string, string) error
	switch {
	case results.IsContainerIssue:
		fixFunc = FixMpegTsToMp4
	case results.IsCodecIssue:
		fixFunc = FixHevcToAvc
	default:
		log.Printf("WARN: Unknown issue type for %s\n", pathToCheck)
		return
	}

	base := strings.TrimSuffix(pathToCheck, filepath.Ext(pathToCheck))
	fixedPath := fmt.Sprintf("%s_fixed_%d.mp4", base, time.Now().Unix())

	if err := fixFunc(pathToCheck, fixedPath); err != nil {
		log.Printf("FATAL: Failed to fix %s: %v\n", pathToCheck, err)
		return
	}

	os.Remove(pathToCheck)
	os.Rename(fixedPath, pathToCheck)
}
