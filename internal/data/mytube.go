package data

type Data struct {
	ModelName  string
	Files      []External
	LocalFiles []External
}

var MyTube Data

type External struct {
	Name     string
	FileUrls string
	Type     string
	XSize    int
	YSize    int
}

var FileData []External

// FileUrls, FileType XSize, YSize
func (d *External) Set(url, fileType string, x, y int) *External {

	d.FileUrls = url

	FileData = append(FileData, *d)
	return d
}
