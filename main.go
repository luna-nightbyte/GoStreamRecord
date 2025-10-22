package main

import (
	// Added: Required for bytes.NewReader (implements io.ReadSeeker)

	"embed"
	_ "embed"
	"fmt"
	"log"
	"remoteCtrl/internal"
	"remoteCtrl/internal/media/localfolder"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/system/prettyprint"
	"remoteCtrl/internal/system/version"
	"remoteCtrl/internal/web"
	"remoteCtrl/internal/web/telegram"
)

//go:embed vue/login/dist/*
var VueLoginFiles embed.FS

//go:embed vue/app/dist/*
var VueDistFiles embed.FS

func init() {

	fmt.Println()
	prettyprint.P.BoldGrey.Println("Software version: ")
	prettyprint.P.Cyan.Println(version.Version)
	prettyprint.P.BoldGrey.Println("Github commit sha:")
	prettyprint.P.Cyan.Println(version.Shasum)
	fmt.Println()
	prettyprint.P.BoldGreen.Println("Startup")
	system.Check_requirements()

}
func main() {
	system.System.WaitForNetwork = false
	err := internal.Init()
	if err != nil {
		log.Fatal(err)
	}
	go localfolder.ContiniousRead(system.System.Config.OutputFolder)
	go web.ServeHTTP(system.System.Context, VueLoginFiles, VueDistFiles)

	<-system.System.Context.Done()

	telegram.Bot.SendMsg("Server shutdown")
}
