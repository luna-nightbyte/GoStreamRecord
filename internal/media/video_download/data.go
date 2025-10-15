package video_download

import (
	"log"
	"os"
	"path/filepath"
)

type video struct {
	MP4pwd string
	Type   string
	XSize  int
	YSize  int
}
type tmpFile struct {
	NoStream bool
	Mp4      video
	Dir           string
	TSContentfile string
	TSSegmentsTXT string
}
type OutputFile struct {
	Type string
	Path string
	Name string
}
 
func (f *tmpFile) RemoveTmps() error {
	err := os.RemoveAll(f.Dir)
	if err != nil {
		return err
	}
	return os.RemoveAll(f.TSSegmentsTXT)

}

func (f *tmpFile) CreateTempDirs() {
	err := os.MkdirAll(f.Dir, os.ModePerm)
	if err != nil {
		log.Println("Error creating output directory:", err)
		return
	}

}

func (f *tmpFile) CreateTempFiles() {
	outputFile, err := os.Create(f.TSContentfile)
	if err != nil {
		log.Println("Error creating f.TSContentfile file:", err)

	}
	defer outputFile.Close()
	outputFile, err = os.Create(f.Dir)
	if err != nil {
		log.Println("Error creating f.TSSegmentfilePath file:", err)
		return
	}
	defer outputFile.Close()
}

func (f *tmpFile) CreateSegment(file string) (*os.File, error) {
	return os.Create(filepath.Join(f.Dir, file))
}
