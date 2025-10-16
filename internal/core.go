package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"remoteCtrl/internal/command"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/system/cookies"
	"remoteCtrl/internal/system/logger"
	"remoteCtrl/internal/system/prettyprint"
	"remoteCtrl/internal/system/settings"
	"remoteCtrl/internal/utils"
	"remoteCtrl/internal/web/telegram"
	"strconv"
	"syscall"
	"time"
)

var onlineCheckIP = "192.168.10.173"
var connected bool

func Init() error {

	system.System.DB = settings.Init()
	// Context for shutdown
	system.System.Context, system.System.Cancel = context.WithCancel(context.Background())

	// COS sig
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("\nReceived a shutdown signal: %v. Initiating graceful shutdown...\n", sig)
		system.System.Cancel()
	}()

	// -- Default init
	if len(os.Args) < 2 {
		isLive := false
		// -- -- Network
		if system.System.WaitForNetwork {
			isLive = utils.Ping(onlineCheckIP)
		}
		attempts := 0
		max := 5
		for !isLive && system.System.WaitForNetwork {
			attempts++
			if attempts == max {
				if telegram.Bot.Enabled() {
					fmt.Println(prettyprint.BoldRed("No network. Telegram disabled!"))
					log.Println("No network. Telegram disabled!")
					system.System.DB.Settings.Telegram.Enabled = false
					telegram.Bot.Disable()
				}
				break
			}
			log.Println("No network connection..")
			time.Sleep(30 * time.Second)
			isLive = utils.Ping(onlineCheckIP)
		}
		cookies.Session = cookies.New(system.System.DB.Settings)
		logger.Init(logger.Log_path)

		// -- -- Telegram
		if system.System.DB.Settings.Telegram.Enabled {
			telegram.Bot.Init()
			telegram.Bot.SendStartup(strconv.Itoa(system.System.DB.Settings.App.Port))
		}
		return nil

	}
	//
	//	Only argument handling below
	//

	cmdName := os.Args[1]
	statup_command, exists := command.CMD.Startup.Map[cmdName]
	if !exists {
		system.StartupError()
		fmt.Println(prettyprint.Cyan("Unknown command:"), cmdName)
		log.Println("Unknown command:", cmdName)
		return fmt.Errorf("unknown command")
	}

	statup_command.Run(os.Args[2:])

	system.System.Cancel()

	return nil
}
