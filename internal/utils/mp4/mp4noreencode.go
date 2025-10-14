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
)

var numberRegex_n = regexp.MustCompile(`(\d+)`)

// orderByNumber orders files by numerical value in their names.
func orderByNumber_n(inputDir string) ([]string, error) {
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
func extractNumber_n(fileName string) int {
	match := numberRegex.FindStringSubmatch(filepath.Base(fileName))
	if len(match) < 2 {
		return -1 // if no number found, return -1 to put it at the beginning
	}
	num, _ := strconv.Atoi(match[1])
	return num
}

// TSToMP4 reads a folder, concatenates all .ts files, and converts the result to an MP4 file
func TSToMP4_n(inputDir, outputFile, tsTXT_tmp string) error {
	// Read the directory
	inputFiles, err := orderByNumber(inputDir)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error reading directory: %v", err)
	}

	if len(inputFiles) == 0 {
		log.Println("no .ts files found in directory")
		return fmt.Errorf("no .ts files found in directory")
	}

	// Create a temporary file to hold the list of input files for ffmpeg
	tempListFile, err := os.Create(tsTXT_tmp)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error creating temporary list file: %v", err)
	}
	defer tempListFile.Close()

	// Write the list of input files to the temporary file
	writer := bufio.NewWriter(tempListFile)
	for _, inputFile := range inputFiles {
		_, err := writer.WriteString(fmt.Sprintf("file '%s'\n", inputFile))
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

	// Concatenate the TS files into a single MP4 file
	ffmpegArgs := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", tsTXT_tmp,
		"-c", "copy",
		outputFile,
		"-y",
	}
	cmd := exec.Command("ffmpeg", ffmpegArgs...)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Println(err)
		return fmt.Errorf("error concatenating files: %v, output: %s", err, string(output))
	}

	log.Println("Done saving!")
	return nil
}
