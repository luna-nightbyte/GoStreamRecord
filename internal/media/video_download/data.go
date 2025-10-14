package video_download

import (
	"log"
	"os"
	"path/filepath"
)

type tmpFiles struct {
	Dir           string
	TSContentfile string
	TSSegmentsTXT string
}
type video struct {
	MP4pwd string
	Type   string
	XSize  int
	YSize  int
}
type File struct {
	NoStream bool
	Mp4      video
	Tmp      tmpFiles
}
type OutputFile struct {
	Type string
	Path string
	Name string
}

var OutputFiles []OutputFile

var (
	outputDir = "videos"
	TMP       File
)

func init() {

	TMP.Tmp.Dir = "tmp"
	TMP.Tmp.TSContentfile = filepath.Join(TMP.Tmp.Dir, "output.ts")
	TMP.Tmp.TSSegmentsTXT = "./tsFiles.txt"

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.Mkdir(outputDir, 0755)
		if err != nil {
			log.Println(err)
		}
	}
	TMP.Tmp.CreateTempDirs()

}

func (f *File) RemoveTmps() error {
	err := os.RemoveAll(f.Tmp.Dir)
	if err != nil {
		return err
	}
	return os.RemoveAll(f.Tmp.TSSegmentsTXT)

}

func (f *tmpFiles) CreateTempDirs() {
	err := os.MkdirAll(f.Dir, os.ModePerm)
	if err != nil {
		log.Println("Error creating output directory:", err)
		return
	}

}

func (f *tmpFiles) CreateTempFiles() {
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

func (f *tmpFiles) CreateSegment(file string) (*os.File, error) {
	return os.Create(filepath.Join(f.Dir, file))
}
