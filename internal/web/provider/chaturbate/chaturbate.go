package chaturbate

import (
	"net/http"
	"strings"
	"time"
)

type Chaturbate struct {
	Url           string `json:"url"`
	StreamURL     string `json:"stream_url"`
	CorrectedName string `json:"username"`
}

var provider_url string = "https://chaturbate.com/"

func (b *Chaturbate) Init(username string) any {
	b.Url = provider_url
	b.StreamURL = provider_url + username
	b.CorrectedName = username
	return b
}

// IsOnline checks if the streamer is online by checking if a thumbnail is available from the stream.
func (c *Chaturbate) IsOnline(username string) bool {
	// Short delay before making the call.
	//Check once if a thumbnail is available
	urlStr := "https://jpeg.live.mmcdn.com/stream?room=" + c.TrueName(username)
	http.DefaultClient.Timeout = time.Duration(1 * time.Second)
	resp, err := http.DefaultClient.Get(urlStr)
	//resp, err := http.Get(urlStr)
	if err != nil {
		//log.Printf("Error in making request: %v", err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK { // Streamer is not online if response if not 200
		return false

	}
	return true
}

// Not necessary for this as of now.
func (b *Chaturbate) TrueName(name string) string {
	return strings.ToLower(name)
}

// Not necessary for this as of now.
func (b *Chaturbate) Settings(provider any) any {
	return b
}
