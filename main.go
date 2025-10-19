package main

import (
	// Added: Required for bytes.NewReader (implements io.ReadSeeker)

	"embed"
	_ "embed"
	"fmt"
	"log"
	"remoteCtrl/internal"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/system/prettyprint"
	"remoteCtrl/internal/system/version"
	"remoteCtrl/internal/utils"
	"remoteCtrl/internal/web"
	"remoteCtrl/internal/web/telegram"
)

//go:embed vue/login/dist/*
var VueLoginFiles embed.FS

//go:embed vue/app/dist/*
var VueDistFiles embed.FS

func init() {

	fmt.Println(prettyprint.Green("Startup"))
	fmt.Println(prettyprint.BoldGrey("Software version", version.Version))
	fmt.Println(prettyprint.BoldGrey("Software sha256:", version.Shasum), "\n")
	ytDLP_path := utils.CheckPath("yt-dlp")

	ffmpeg_path := utils.CheckPath("ffmpeg")
	if ytDLP_path == "" {
		fmt.Println("missing yt-dlp. Please instalwebControllerl before running this program.")
	}
	if ffmpeg_path == "" {
		fmt.Println("missing yt-dlp. Please install before running this program.")
	}
	if ytDLP_path == "" || ffmpeg_path == "" {
		log.Fatal("missing dependencies")
	}
}
func main() {
	system.System.WaitForNetwork = false
	err := internal.Init()
	if err != nil {
		log.Fatal(err)
	}

	go web.ServeHTTP(system.System.Context, VueLoginFiles, VueDistFiles)

	<-system.System.Context.Done()

	telegram.Bot.SendMsg("Server shutdown")
}
