package settings

// Version 1.
type Settings struct {
	App app `json:"app"`
}

type app struct {
	Port          int        `json:"port"`
	Loop_interval int        `json:"loop_interval_in_minutes"`
	Videos_folder string     `json:"video_output_folder"`
	RateLimit     rate_limit `json:"rate_limit"`
	UseEnvKey     bool       `json:"use_env_key"`
}
type rate_limit struct {
	Enable bool `json:"enable"`
	Time   int  `json:"time"`
}
