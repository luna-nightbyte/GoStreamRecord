package settings

// version 2.
type Settings_v2 struct {
	App web_ui `json:"app"`
}

type web_ui struct {
	Port          int    `json:"port"`
	Videos_folder string `json:"output_folder"`
}
