package online

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Online struct {
	mu           sync.Mutex
	Enabled      bool
	service      *drive.Service
	lastModified map[string]time.Time
}

type Settings struct {
	Enable bool
}

func (o *Online) New(credentialsPath string) error {
	ctx := context.Background()
	b, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		return err
	}

	o.service, err = drive.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		return err
	}
	o.lastModified = make(map[string]time.Time)
	return nil
}
func (o *Online) ListFilesInFolder(parentId string) ([]string, error) {
	q := fmt.Sprintf("'%s' in parents and trashed=false", parentId)
	r, err := o.service.Files.List().Q(q).Fields("files(id, name)").Do()
	if err != nil {
		return nil, err
	}
	var files []string
	for _, file := range r.Files {
		files = append(files, file.Name)
	}
	return files, nil
}

func (o *Online) GetLastModified(fileId string) (time.Time, error) {
	f, err := o.service.Files.Get(fileId).Fields("modifiedTime").Do()
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.Parse(time.RFC3339, f.ModifiedTime)
	if err != nil {
		return time.Time{}, err
	}
	o.lastModified[fileId] = t
	return t, nil
}

func (o *Online) CreateFolderIfNotExist(name string, parentId string) (string, error) {
	q := fmt.Sprintf("name='%s' and mimeType='application/vnd.google-apps.folder' and trashed=false", name)
	if parentId != "" {
		q += fmt.Sprintf(" and '%s' in parents", parentId)
	}

	r, err := o.service.Files.List().Q(q).Do()
	if err != nil {
		return "", err
	}

	if len(r.Files) > 0 {
		return r.Files[0].Id, nil
	}

	folder := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentId},
	}

	created, err := o.service.Files.Create(folder).Do()
	if err != nil {
		return "", err
	}

	return created.Id, nil
}

func (o *Online) FindFileByName(name, parentId string) (string, error) {
	q := fmt.Sprintf("name='%s' and trashed=false", name)
	if parentId != "" {
		q += fmt.Sprintf(" and '%s' in parents", parentId)
	}

	r, err := o.service.Files.List().Q(q).Fields("files(id, name)").Do()
	if err != nil {
		return "", err
	}

	if len(r.Files) == 0 {
		return "", errors.New("file not found")
	}

	return r.Files[0].Id, nil
}

func (o *Online) DownloadFile(fileId, destPath string) error {
	res, err := o.service.Files.Get(fileId).Download()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	return err
}

func (o *Online) ToggleEnable() {
	o.Enabled = !o.Enabled
}

func (d *Online) UploadFile(localPath, folderName string) (string, error) {
	file, err := os.Open(localPath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	folderID, err := d.CreateFolderIfNotExist(folderName, "")
	if err != nil {
		return "", err
	}
	f := &drive.File{
		Name:    filepath.Base(localPath),
		Parents: []string{folderID},
	}

	uploaded, err := d.service.Files.Create(f).Media(file).Do()
	if err != nil {
		return "", err
	}

	return uploaded.Id, nil
}
func (d *Online) WriteJSON(fileId string, data []byte) error {
	_, err := d.service.Files.Update(fileId, nil).Media(bytes.NewReader(data)).Do()
	return err
}

func (d *Online) ReadJSON(fileId string) (*[]byte, error) {
	res, err := d.service.Files.Get(fileId).Download()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var data []byte
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
