package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"remoteCtrl/internal/system"
	"runtime"
)

var (
	logFile  *os.File
	Log_path string = "./app.log"
)

type logWriter struct{}

func Init(logPath string) {
	var err error
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.SetOutput(logWriter{})
	log.SetFlags(0)
}

func (w logWriter) Write(p []byte) (n int, err error) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		// line = 0
	} else {
		file = trimPath(file)
	}
	if system.DEBUG {
		formattedMsg := fmt.Sprintf("\"./%s\":[%d] %s", file, line, p)
		return logFile.Write([]byte(formattedMsg))
	}
	// formattedMsg := fmt.Sprintf("\"./%s\":[%d] %s", file, line, p)
	formattedMsg := fmt.Sprintf("%s", p)
	return logFile.Write([]byte(formattedMsg))
}

func trimPath(fullPath string) string {
	wd, err := os.Getwd()
	if err != nil {
		return fullPath
	}
	relPath, err := filepath.Rel(wd, fullPath)
	if err != nil {
		return fullPath
	}
	return relPath
}
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
