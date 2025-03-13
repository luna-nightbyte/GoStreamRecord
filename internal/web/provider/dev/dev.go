package dev

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Dev struct {
	Url        string `json:"url"`
	StreamURL  string `json:"stream_url"`
	StreamName string `json:"username"`
}

var provider_url string = "http://localhost:8050/"
var stream_url string = "http://localhost:8050/video_feed"

// Struct matching the JSON structure (partial example, you can expand this as needed)
type InitialState struct {
	ChatLocalData struct {
		IsAway bool   `json:"isAway"`
		Gender string `json:"gender"`
	} `json:"chatLocalData"`
	ChatHost struct {
		Username string `json:"username"`
		Online   bool   `json:"online"`
		ShowType string `json:"showType"`
	} `json:"chatHost"`
}

func (b *Dev) Init(username string) any {
	b.Url = provider_url
	b.StreamURL = stream_url
	b.StreamName = b.TrueName(username)
	return b
}

// IsOnline checks if the streamer is online by checking if a thumbnail is available from the stream.
func (b *Dev) IsOnline(username string) bool {
	state := b.getState(username)
	if (state == InitialState{}) {
		log.Println("Error recieving online state.")
		return false
	}
	b.StreamName = state.ChatHost.Username
	return state.ChatHost.Online && state.ChatHost.ShowType == "public" && !state.ChatLocalData.IsAway
}

func (b *Dev) getState(username string) InitialState {
	resp, err := http.Get(fmt.Sprintf("%s%s", provider_url, strings.ToLower(username)))
	if err != nil {
		log.Printf("HTTP request failed: %v\n", err)
		return InitialState{}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading: %v\n", err)
		return InitialState{}
	}

	var state InitialState
	if err := json.Unmarshal(bodyBytes, &state); err != nil {
		log.Println(fmt.Errorf("JSON unmarshal failed: %v", err))
		return InitialState{}
	}
	return state
}

func (b *Dev) TrueName(name string) string {
	if name != b.StreamName {
		b.StreamName = strings.ToLower(name)
	}
	return b.StreamName
}

// Not necessary for this as of now.
func (b *Dev) Settings(provider any) any {
	return b
}
