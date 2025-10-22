package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"remoteCtrl/internal/command"
	"remoteCtrl/internal/db"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/system/cookies"
	"remoteCtrl/internal/system/logger"
	"remoteCtrl/internal/system/prettyprint"
	"remoteCtrl/internal/utils"
	"remoteCtrl/internal/web/handlers/status"
	"remoteCtrl/internal/web/telegram"

	"strconv"
	"syscall"
	"time"
)

var onlineCheckIP = "8.8.8.8"

func Init() error {

	//system.System.Config = settings.Init()
	// Context for shutdown
	system.System.Context, system.System.Cancel = context.WithCancel(context.Background())

	db.Init(system.System.Context, "")

	cfg, err := db.DataBase.Config()
	if err != nil {
		log.Fatal(err)
	}
	system.System.Config = cfg
	// db.AEAKEY:= getSecret(secretKey)
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
		// -- -- Network
		ticker := time.NewTicker(30 * time.Second)
		go func() {
			status.Status.IsOnline = utils.Ping(onlineCheckIP)
			for range ticker.C {
				status.Status.IsOnline = utils.Ping(onlineCheckIP)
			}
		}()

		cookies.Session = cookies.New(system.System.Config)
		logger.Init(logger.Log_path)

		// -- -- Telegram
		if system.System.Config.EnableTelegram {
			telegram.Bot.Init()
			telegram.Bot.SendStartup(strconv.Itoa(system.System.Config.Port))
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
