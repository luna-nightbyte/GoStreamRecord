package file

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	logFile *os.File
)

// logWriter implements io.Writer to format logs with timestamp & caller info
type logWriter struct{}

func (w logWriter) Write(p []byte) (n int, err error) {
	_, file, line, ok := runtime.Caller(3) // Adjust stack depth to get the actual caller
	if !ok {
		file = "???"
		line = 0
	} else {
		file = trimPath(file) // Convert absolute to relative path
	}

	formattedMsg := fmt.Sprintf("\"./%s\":[%d] %s", file, line, p)
	return logFile.Write([]byte(formattedMsg))
}

func trimPath(fullPath string) string {
	wd, err := os.Getwd() // Get working directory
	if err != nil {
		return fullPath // Return absolute path if error occurs
	}
	relPath, err := filepath.Rel(wd, fullPath)
	if err != nil {
		return fullPath // Return absolute path if relative conversion fails
	}
	return relPath
}

func InitLog(logPath string) {
	var err error
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Set our custom writer as the output for Go's log package
	log.SetOutput(logWriter{})
	log.SetFlags(0) // Disable default flags, since we're adding our own timestamp
}

func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
