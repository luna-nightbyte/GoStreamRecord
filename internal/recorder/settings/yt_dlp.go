package recorder_settings

import (
	"GoStreamRecord/internal/db"
	"log"
)

var YT_DLP_Args = map[string]Arguments{}

const (
	NO_PART               = "--no-part"
	MIN_FILESIZE          = "--min-filesize"
	MAX_FILESIZE          = "--max-filesize"
	COOKIES               = "--cookies"
	NO_COOKIES            = "--no-cookies"
	COOKIES_FROM_BRROWSER = "--cookies-from-browser"
	WRITE_THUMBNAIL       = "--write-thumbnail"
	NO_WRITE_THUMBNAIL    = "--no-write-thumbnail"
	WRITE_ALL_THUMBNAILS  = "--write-all-thumbnails"
	VIDEO_MULTISTREAMS    = "--video-multistreams"
	NO_VIDEO_MULTISTREAMS = "--no-video-multistreams"
	AUDIO_MULTISTREAMS    = "--audio-multistreams"
	NO_AUDIO_MULTISTREAMS = "--no-audio-multistreams"
	EMBED_THUMBNAIL       = "--embed-thumbnail"
	NO_EMBED_THUMBNAIL    = "--no-embed-thumbnail"
)

type Arguments struct {
	Arg         string `json:"arg"`
	Usage       string `json:"usage"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

func init() {
	err := db.Read("settings", "yt-dlp.json", &YT_DLP_Args)
	if err != nil {
		log.Fatal(err)
	}
}
