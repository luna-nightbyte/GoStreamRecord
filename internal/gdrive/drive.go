package gdrive

import (
	"fmt"
	"os"
	"remoteCtrl/internal/gdrive/mounted"
	"remoteCtrl/internal/gdrive/online"
	"remoteCtrl/internal/system"
	"sync"
	"time"
)

type FileInfo struct {
	ID           string
	Name         string
	ModifiedTime time.Time
}
type iDrive struct {
	d Drive
}
type Drive interface {
	ListFilesInFolder(parentFolderID string) ([]string, error)
	GetLastModified(fileID string) (time.Time, error)
	CreateFolderIfNotExist(name, parentFolderID string) (string, error)
	FindFileByName(name, parentFolderID string) (string, error)
	DownloadFile(fileID, destPath string) error
	ToggleEnable()
	New(string) error
	ReadJSON(fileId string) (*[]byte, error)                     //
	WriteJSON(fileId string, data []byte) error                  //
	UploadFile(localPath, parentFolderId string) (string, error) //
}

const RootFolder string = "GoStreaRecord"

var (
	CredentialsFile = "settings/credentials.json"

	// Service is the selected implementation (online or mounted)
	Service iDrive

	makeMu sync.Mutex
)

// Init picks online if credentials.json exists; otherwise use mounted.
// Adjust the mounted root if needed.
func SetType(serviceType string) error {
	switch serviceType {
	case "online":
		if _, err := os.Stat(CredentialsFile); err == nil {
			Service.d = new(online.Online)
			err := Service.d.New(CredentialsFile)
			if err != nil {
				return fmt.Errorf("failed to init online drive: %w", err)
			}
		} else {
			return fmt.Errorf("No credentials file found. Google drive disabled. Save new credentials to: %s", CredentialsFile)
		}
	case "mounted":
		Service.d = new(mounted.Mounted)
		err := Service.d.New(system.System.DB.Settings.GoogleDrive.Filepath)
		if err != nil {
			return fmt.Errorf("failed to init mounted drive: %w", err)
		}
	default:
		return fmt.Errorf("not a valid type")
	}
	return nil
}

func (d *iDrive) ListFilesInFolder(parentFolderID string) ([]string, error) {
	return d.d.ListFilesInFolder(parentFolderID)
}
func (d *iDrive) GetLastModified(fileID string) (time.Time, error) {
	return d.d.GetLastModified(fileID)
}
func (d *iDrive) CreateFolderIfNotExist(name, parentFolderID string) (string, error) {
	return d.d.CreateFolderIfNotExist(name, parentFolderID)
}
func (d *iDrive) FindFileByName(name, parentFolderID string) (string, error) {
	return d.d.FindFileByName(name, parentFolderID)
}
func (d *iDrive) DownloadFile(fileID, destPath string) error {
	return d.d.DownloadFile(fileID, destPath)
}
func (d *iDrive) ToggleEnable() {
	d.d.ToggleEnable()
}
func (d *iDrive) UploadFile(localPath, parentFolderId string) (string, error) {
	file, err := os.Open(localPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return Service.d.UploadFile(localPath, parentFolderId)
}
