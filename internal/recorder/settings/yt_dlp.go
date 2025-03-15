package recorder_settings

// argument list for YT-dlp. 
// To add a new argument, update condig and add the argument (case sensitive) to this list. 
const (
	YTDLP_NO_PART               = "--no-part"
	YTDLP_MIN_FILESIZE          = "--min-filesize"
	YTDLP_MAX_FILESIZE          = "--max-filesize"
	YTDLP_COOKIES               = "--cookies"
	YTDLP_NO_COOKIES            = "--no-cookies"
	YTDLP_COOKIES_FROM_BRROWSER = "--cookies-from-browser"
	YTDLP_WRITE_THUMBNAIL       = "--write-thumbnail"
	YTDLP_NO_WRITE_THUMBNAIL    = "--no-write-thumbnail"
	YTDLP_WRITE_ALL_THUMBNAILS  = "--write-all-thumbnails"
	YTDLP_VIDEO_MULTISTREAMS    = "--video-multistreams"
	YTDLP_NO_VIDEO_MULTISTREAMS = "--no-video-multistreams"
	YTDLP_AUDIO_MULTISTREAMS    = "--audio-multistreams"
	YTDLP_NO_AUDIO_MULTISTREAMS = "--no-audio-multistreams"
	YTDLP_EMBED_THUMBNAIL       = "--embed-thumbnail"
	YTDLP_NO_EMBED_THUMBNAIL    = "--no-embed-thumbnail"
)
