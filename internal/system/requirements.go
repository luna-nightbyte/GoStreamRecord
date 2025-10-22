package system

import (
	"fmt"
	"os"
	"remoteCtrl/internal/system/prettyprint"
	"remoteCtrl/internal/utils"
)

var requirements = make(map[string]bool)

const (
	ytDLP   = "yt-dlp"
	ffmpeg  = "ffmpeg"
	ffprobe = "ffprobe"
	curl    = "curl"
)

func Check_requirements() {
	requirements[ytDLP] = utils.Is_installed(ytDLP)
	requirements[ffmpeg] = utils.Is_installed(ffmpeg)
	requirements[ffprobe] = utils.Is_installed(ffprobe)
	requirements[curl] = utils.Is_installed(curl)
	depOK := true
	for rec, ok := range requirements {
		if !ok {
			if depOK {
				prettyprint.P.BoldRed.Println("Dependencies error!")
			}
			depOK = false
			fmt.Printf("	- missing \"%s\". Please install before running this program\n", rec)
		}
	}
	if !depOK {
		prettyprint.P.BoldRed.Println("missing dependencies")
		os.Exit(0)
	}
}
