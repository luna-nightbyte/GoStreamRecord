package mp4

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type FFProbeOutput struct {
	Streams []struct {
		CodecType string `json:"codec_type"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
	} `json:"streams"`
}

func GetVideoResolution(filePath string) (int, int, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "json", filePath)
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to execute ffprobe: %v", err)
	}

	var ffprobeOutput FFProbeOutput
	if err := json.Unmarshal(output, &ffprobeOutput); err != nil {
		return 0, 0, fmt.Errorf("failed to parse ffprobe output: %v", err)
	}

	if len(ffprobeOutput.Streams) == 0 {
		return 0, 0, fmt.Errorf("no video stream found")
	}

	return ffprobeOutput.Streams[0].Width, ffprobeOutput.Streams[0].Height, nil
}
