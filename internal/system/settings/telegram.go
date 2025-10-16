package settings

type telegram struct {
	ChatID  string `json:"chatID"`
	Token   string `json:"token"`
	Enabled bool   `json:"enabled"`
}
