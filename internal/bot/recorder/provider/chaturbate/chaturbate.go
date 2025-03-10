package chaturbate

import (
	"GoStreamRecord/internal/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Chaturbate struct {
	Url           string `json:"url"`
	CorrectedName string `json:"username"`
}

var provider_url string = "https://chaturbate.com/"


func (b *Chaturbate) Init(username string) any {
	b.Url = provider_url
	b.CorrectedName = username
	return b
}

// IsOnline checks if the streamer is online by checking if a thumbnail is available from the stream.
func (c *Chaturbate) IsOnline(username string) bool {
	// Short delay before making the call.

	//Check once if a thumbnail is available
	urlStr := "https://jpeg.live.mmcdn.com/stream?room=" + c.TrueName(username)
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Printf("Error in making request: %v", err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusOK { // Streamer is not online if response if not 200
		return false

	}
	return true
}

// Old method from recordurbate. Not really used in this app.
// IsRoomPublic checks if a given room is public by sending a POST request.
func (c *Chaturbate) IsRoomPublic(username string) bool {
	// Wait for the dbured rate limit.
	time.Sleep(time.Duration(db.Config.Settings.App.RateLimit.Time) * time.Second)
	fmt.Println("Checking ", username)
	urlStr := "https://chaturbate.com/get_edge_hls_url_ajax/"
	data := url.Values{}
	data.Set("room_slug", username)
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return false
	}
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making the request: %v", err)
		return false
	}
	defer resp.Body.Close()

	var res struct {
		Success    bool   `json:"success"`
		RoomStatus string `json:"room_status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return false
	}
	return res.Success && res.RoomStatus == "public"
}

// Not necessary for this as of now.
func (b *Chaturbate) TrueName(name string) string {
	return strings.ToLower(name)
}

// Not necessary for this as of now.
func (b *Chaturbate) Settings(provider any) any {
	return b
}
