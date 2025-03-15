package settings

import (
	"encoding/json"
	"os"
)

type YP_DLP struct {
	YT_DLP map[string]Arguments `json:"yt_dlp"`
}

type Arguments struct {
	Arg         string `json:"arg"`
	Usage       string `json:"usage"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

func (yt *YP_DLP) Read(filepath string, i interface{}) error {
	contentBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(contentBytes, i)
	if err != nil {
		return err
	}
	return nil
}
