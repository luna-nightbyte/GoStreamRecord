package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/utils"
	"time"
)

// Telegram token, chat(s) and API
const (
	api_root        string = "https://api.telegram.org/bot"
	api_sendPhoto   string = "sendPhoto"
	api_sendMessage string = "sendMessage"
)

var Bot bot

// Default values
const (
	lowSpaceWarningThresholdGB = 1500
	httpTimeout                = 1 * time.Minute
)

type bot struct {
	enabled bool
	IP      string
	Message map[string]interface{}
}

func newMsg() map[string]interface{} {
	return map[string]interface{}{
		"chat_id":    system.System.Config.TelegramChatID,
		"parse_mode": "Markdown",
	}
}

func (b *bot) Init() {
	b.enabled = true
	b.IP = utils.GetLocalIp()
}

func (b *bot) Enabled() bool {
	return b.enabled
}
func (b *bot) Disable() {
	b.enabled = false
}

// SendTelegramStartup sends a startup notification
func (b bot) SendStartup(port string) {
	if !b.enabled {
		return
	}
	message := fmt.Sprintf("GoStreamRecord startup\nVisit your local webserver at http://%s:%s", b.IP, port)

	go b.sendMessage(message)
}

// SendTelegramStartup sends a startup notification
func (b bot) SendMsg(msg string) {
	if !b.enabled {
		return
	}
	message := fmt.Sprintf("🟢 *Message*\n[%s]:\n\n%s", b.IP, msg)
	go b.sendMessage(message)
}

func createApiRequest(ctx context.Context, payload map[string]interface{}, endpoint string) *http.Request {

	url := fmt.Sprintf("%s%s/%s", api_root, system.System.Config.TelegramToken, endpoint)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling Telegram payload: %v", err)
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating Telegram request: %v", err)
		return nil
	}
	return req
}

// sendTelegramMessage sends any message to Telegram with timeout
func (b *bot) sendMessage(message string) {
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()
	payload := newMsg()
	payload["text"] = message
	req := createApiRequest(ctx, payload, "sendMessage")
	req.Header.Set("Content-Type", "application/json")

	payload["text"] = ""
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending Telegram notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Telegram API returned status: %d", resp.StatusCode)
		return
	}

	log.Println("Telegram message sent successfully")
}

// sendTelegramMessage sends any message to Telegram with timeout
func (b bot) SendPhoto(imagePath, caption string) {
	if !b.enabled {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()
	payload := newMsg()
	payload["photo"] = imagePath

	if caption != "" {
		payload["caption"] = caption
	}
	req := createApiRequest(ctx, payload, "sendPhoto")
	req.Header.Set("Content-Type", "application/json")

	payload["caption"] = ""

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending Telegram notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Telegram API returned status: %d", resp.StatusCode)
		return
	}

	log.Println("Telegram message sent successfully")
}

func (b *bot) sendVideo(videoPath, caption string) {
	if !b.enabled {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()
	payload := newMsg()
	payload["photo"] = videoPath

	if caption != "" {
		payload["caption"] = caption
	}

	req := createApiRequest(ctx, payload, "sendVideo")
	req.Header.Set("Content-Type", "application/json")

	payload["caption"] = ""

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending Telegram notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Telegram API returned status: %d", resp.StatusCode)
		return
	}

	log.Println("Telegram message sent successfully")
}
