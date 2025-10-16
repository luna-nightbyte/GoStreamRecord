package mounted

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Mounted struct {
	mu           sync.Mutex
	Enabled      bool
	lastModified map[string]time.Time
}

func (o *Mounted) New(NoteUsed string) error {
	o = &Mounted{lastModified: make(map[string]time.Time)}
	return nil
}
func (o *Mounted) ToggleEnable() {
	o.Enabled = !o.Enabled
}
func (o *Mounted) ListFilesInFolder(folderPath string) ([]string, error) {
	folderContent, err := os.ReadDir(folderPath)
	if err != nil {
		return []string{}, err
	}
	var files []string
	for _, dirEntry := range folderContent {
		if dirEntry.IsDir() {
			dirFiles, _ := o.ListFilesInFolder(filepath.Join(folderPath, dirEntry.Name()))
			files = append(files, dirFiles...)
			continue
		}
		files = append(files, filepath.Join(folderPath, dirEntry.Name()))
	}
	return files, nil
}

func (o *Mounted) GetLastModified(fileId string) (time.Time, error) {

	return time.Time{}, nil
}

func (o *Mounted) CreateFolderIfNotExist(name string, parentId string) (string, error) {
	_, err := os.Stat(parentId)
	if os.IsNotExist(err) {
		return "", nil
	}
	return "", os.MkdirAll(parentId, 0755)
}

func (o *Mounted) FindFileByName(name, parentId string) (string, error) {
	files, err := o.ListFilesInFolder(parentId)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if strings.Contains(file, name) {
			return file, nil
		}
	}
	return "", nil
}

func (o *Mounted) DownloadFile(fileId, destPath string) error {
	// Open the source file for reading.
	sourceFile, err := os.Open(fileId)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create the destination file for writing.
	destinationFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("could not create destination file: %w", err)
	}
	defer destinationFile.Close()

	// Use io.Copy to efficiently copy the file content.
	if _, err := io.Copy(destinationFile, sourceFile); err != nil {
		return fmt.Errorf("could not copy file content: %w", err)
	}

	// Flush and close the destination file to ensure all data is written.
	if err := destinationFile.Close(); err != nil {
		return fmt.Errorf("could not close destination file: %w", err)
	}

	return err
}
func (o *Mounted) ReadJSON(filePath string) (*[]byte, error) {
	if !strings.Contains(filepath.Base(filePath), ".json") {
		return nil, fmt.Errorf("not a json file..")
	}
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return &jsonData, nil
}
func (o *Mounted) WriteJSON(filePath string, data []byte) error {
	var jsonFile *os.File
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		jsonFile = file
	}
	_, err = jsonFile.Write(data)
	return err
}

// Uploadfile name is Loced, in Mounted that means copy to mounted Dir.
func (o *Mounted) UploadFile(localPath, destinationFolder string) (string, error) {
	mountedPath := filepath.Join(destinationFolder, filepath.Base(localPath))

	_, err := os.Stat(localPath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("local path does not exist")
	}
	_, err = os.Stat(mountedPath)
	if !os.IsNotExist(err) {
		os.Remove(mountedPath)
	}
	fileContent, err := os.ReadFile(localPath)
	if err != nil {
		return "", err
	}
	file, err := os.Create(mountedPath)
	if err != nil {
		return "", err
	}

	_, err = file.Write(fileContent)
	return "", err
}
