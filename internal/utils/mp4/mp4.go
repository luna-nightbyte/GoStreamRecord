package mp4

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var numberRegex = regexp.MustCompile(`(\d+)`)

// orderByNumber orders files by numerical value in their names.
func orderByNumber(inputDir string) ([]string, error) {
	files, err := os.ReadDir(inputDir)
	if err != nil {
		return nil, err
	}

	var tsFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".ts") {
			tsFiles = append(tsFiles, filepath.Join(inputDir, file.Name()))
		}
	}

	// Sort files by their numeric value extracted from the filename
	sort.Slice(tsFiles, func(i, j int) bool {
		return extractNumber(tsFiles[i]) < extractNumber(tsFiles[j])
	})

	return tsFiles, nil
}

// extractNumber extracts the first number from the filename using a regex.
func extractNumber(fileName string) int {
	match := numberRegex.FindStringSubmatch(filepath.Base(fileName))
	if len(match) < 2 {
		return -1 // if no number found, return -1 to put it at the beginning
	}
	num, _ := strconv.Atoi(match[1])
	return num
}

// Re-encode TS files to fix timestamps and other issues.
func reencodeTSFiles(inputFiles []string, inputDir string, concurrency int) ([]string, error) {
	var reencodedFiles []string
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	sem := make(chan struct{}, concurrency)
	progressCh := make(chan int, len(inputFiles))
	errorsCh := make(chan error, len(inputFiles))

	for i, inputFile := range inputFiles {
		wg.Add(1)
		sem <- struct{}{} // Acquire a token

		go func(i int, inputFile string) {
			defer wg.Done()
			defer func() { <-sem }() // Release the token

			startTime := time.Now()
			s := strings.Split(inputFile, "/")
			outputFile := fmt.Sprintf("%s/fixed_%s", inputDir, s[len(s)-1])
			ffmpegArgs := []string{
				"-i", inputFile,
				"-c:v", "libx264",
				"-c:a", "aac",
				"-strict", "experimental",
				"-b:a", "192k",
				outputFile,
				"-y",
			}
			cmd := exec.Command("ffmpeg", ffmpegArgs...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				errorsCh <- fmt.Errorf("error re-encoding file %s: %v, output: %s", inputFile, err, string(output))
				return
			}
			elapsedTime := time.Since(startTime).Seconds()
			log.Printf("Re-encoded %s to %s in %.2f seconds\n", inputFile, outputFile, elapsedTime)
			progressCh <- i

			mu.Lock()
			reencodedFiles = append(reencodedFiles, outputFile)
			mu.Unlock()
		}(i, inputFile)
	}

	go func() {
		wg.Wait()
		close(progressCh)
		close(errorsCh)
	}()

	totalFiles := len(inputFiles)
	for {
		select {
		case _, ok := <-progressCh:
			if !ok {
				progressCh = nil
			} else {
				progress := float64(len(reencodedFiles)) / float64(totalFiles) * 100
				fmt.Printf("Progress: %.2f%%\n", progress)
			}
		case err, ok := <-errorsCh:
			if !ok {
				errorsCh = nil
			} else if err != nil {
				return nil, err
			}
		}
		if progressCh == nil && errorsCh == nil {
			break
		}
	}

	// Sort the reencoded files to ensure they are in the correct order
	sort.Slice(reencodedFiles, func(i, j int) bool {
		return extractNumber(reencodedFiles[i]) < extractNumber(reencodedFiles[j])
	})

	return reencodedFiles, nil
}

// Monitor FFmpeg progress.
func monitorProgress(cmd *exec.Cmd) error {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stderr)
	durationRegex := regexp.MustCompile(`Duration: (\d+):(\d+):(\d+.\d+)`)
	timeRegex := regexp.MustCompile(`time=(\d+):(\d+):(\d+.\d+)`)

	var duration float64
	for scanner.Scan() {
		line := scanner.Text()
		if matches := durationRegex.FindStringSubmatch(line); matches != nil {
			hours := matches[1]
			minutes := matches[2]
			seconds := matches[3]
			duration = parseTime(hours, minutes, seconds)
		} else if matches := timeRegex.FindStringSubmatch(line); matches != nil {
			hours := matches[1]
			minutes := matches[2]
			seconds := matches[3]
			currentTime := parseTime(hours, minutes, seconds)
			if duration > 0 {
				p := (currentTime / duration) * 100
				fmt.Printf("Progress: %.2f%%\n", p)
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

// parseTime converts a time string to seconds.
func parseTime(hours, minutes, seconds string) float64 {
	var h, m int
	var s float64
	fmt.Sscanf(hours, "%d", &h)
	fmt.Sscanf(minutes, "%d", &m)
	fmt.Sscanf(seconds, "%f", &s)
	return float64(h)*3600 + float64(m)*60 + s
}

// TSToMP4 reads a folder, re-encodes all .ts files, concatenates them, and converts the result to an MP4 file
func TSToMP4(inputDir, outputFile, tsTXT_tmp string, concurrency int) error {
	// Read the directory
	inputFiles, err := orderByNumber(inputDir)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error reading directory: %v", err)
	}

	if len(inputFiles) == 0 {
		log.Println(err)
		return fmt.Errorf("no .ts files found in directory")
	}

	// Re-encode TS files
	log.Println(inputDir)
	reencodedFiles, err := reencodeTSFiles(inputFiles, inputDir, concurrency)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error re-encoding TS files: %v", err)
	}

	// Create a temporary file to hold the list of re-encoded input files for ffmpeg
	tempListFile, err := os.Create(tsTXT_tmp)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error creating temporary list file: %v", err)
	}
	defer tempListFile.Close()

	// Write the list of re-encoded input files to the temporary file
	writer := bufio.NewWriter(tempListFile)
	for _, reencodedFile := range reencodedFiles {
		_, err := writer.WriteString(fmt.Sprintf("file '%s'\n", reencodedFile))
		if err != nil {
			log.Println(err)
			return fmt.Errorf("error writing to temporary list file: %v", err)
		}
	}

	// Close the writer to ensure it's written to disk
	if err := writer.Flush(); err != nil {
		log.Println(err)
		return fmt.Errorf("error flushing temporary list file: %v", err)
	}

	// Concatenate the re-encoded TS files into a single MP4 file
	ffmpegArgs := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", tsTXT_tmp,
		"-c", "copy",
		outputFile,
		"-y",
	}
	cmd := exec.Command("ffmpeg", ffmpegArgs...)
	if err := monitorProgress(cmd); err != nil {
		log.Println(err)
		return fmt.Errorf("error concatenating files: %v", err)
	}

	log.Println("Done saving!")
	return nil
}
