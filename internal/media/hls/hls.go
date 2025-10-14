package hls

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Downloader represents an HLS stream downloader.
type Downloader struct {
	URL     string // HLS playlist (m3u8)
	Output  string // output file path (e.g., "output.mp4")
	Timeout time.Duration
}

// New creates a new Downloader with sane defaults.
func New(url, output string) *Downloader {
	if output == "" {
		output = "output.mp4"
	}
	return &Downloader{
		URL:     url,
		Output:  output,
		Timeout: 0, // no timeout by default
	}
}

// Start begins downloading the HLS stream into the output file.
// It uses ffmpeg to handle the HLS transport and mux audio/video.
func (d *Downloader) Start() error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH: %w", err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(d.Output), 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	args := []string{
		"-y",        // overwrite existing file
		"-i", d.URL, // input HLS stream
		"-c", "copy", // copy codecs (no re-encoding)
		"-bsf:a", "aac_adtstoasc", // fix audio format if needed
		d.Output,
	}

	ctx := context.Background()
	if d.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, d.Timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Downloading stream: %s → %s\n", d.URL, d.Output)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w", err)
	}

	fmt.Println("✅ Download complete:", d.Output)
	return nil
}
