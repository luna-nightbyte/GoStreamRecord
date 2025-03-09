package bongacams

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type BongaCams struct {
	Url           string `json:"url"`
	CorrectedName string `json:"username"`
}

var provider_url string = "https://bongacams.com/"

func (b *BongaCams) Init( username string) any {
	b.Url = provider_url
	b.CorrectedName = b.TrueName(username)
	return b
}

// IsOnline checks if the streamer is online by checking if a thumbnail is available from the stream.
func (b *BongaCams) IsOnline(username string) bool {
	state := b.getState(username)
	fmt.Println(state)
	b.CorrectedName = state.ChatHost.Username
	return state.ChatHost.Online && state.ChatHost.ShowType == "public"
}

// Struct matching the JSON structure (partial example, you can expand this as needed)
type InitialState struct {
	ChatLocalData struct {
		IsAway    bool   `json:"isAway"`
		IsOffline bool   `json:"isOffline"`
		Gender    string `json:"gender"`
	} `json:"chatLocalData"`
	ChatHost struct {
		Username    string `json:"username"`
		DisplayName string `json:"displayName"`
		Quality     string `json:"quality"`
		Online      bool   `json:"online"`
		ShowType    string `json:"showType"`
		LoversCount int    `json:"loversCount"`
		Gender      string `json:"gender"`
	} `json:"chatHost"`
}

func (b *BongaCams) getState(username string) InitialState {
	resp, err := http.Get(fmt.Sprintf("https://no.bongacams.com/%s", strings.ToLower(username)))
	if err != nil {
		log.Printf("HTTP request failed: %v\n", err)
		return InitialState{}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return InitialState{}
	}
	body := string(bodyBytes)

	// Regular expression to extract JSON from <script> tag
	re := regexp.MustCompile(`<script[^>]*data-type="initialState"[^>]*type="application/json"[^>]*>(.*?)</script>`)
	matches := re.FindStringSubmatch(body)
	if len(matches) < 2 {
		log.Println(fmt.Errorf("Initial state JSON not found in HTML"))
		return InitialState{}
	}
	jsonData := matches[1]

	var state InitialState
	if err := json.Unmarshal([]byte(jsonData), &state); err != nil {
		log.Println(fmt.Errorf("JSON unmarshal failed: %v", err))
		return InitialState{}
	}
	return state
}

func (b *BongaCams) TrueName(name string) string {
	if name != b.CorrectedName {
		new_state := b.getState(name)
		b.CorrectedName = new_state.ChatHost.Username
	}
	return b.CorrectedName
}

// Not necessary for this as of now.
func (b *BongaCams) Settings(provider any) any {
	return b
}
