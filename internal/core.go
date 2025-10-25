package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"remoteCtrl/internal/command"
	"remoteCtrl/internal/db"
	"remoteCtrl/internal/media/localfolder"
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

	logger.Init(logger.Log_path)
	// Context for shutdown
	system.System.Context, system.System.Cancel = context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("\nReceived a shutdown signal: %v. Initiating graceful shutdown...\n", sig)
		system.System.Cancel()
	}()

	// Database
	db.Init(system.System.Context, "")
	cfg := db.DataBase.Config()
	if cfg.Port == -1 {
		log.Fatal("Config load failed")
	}
	system.System.Config = cfg

	// -- Backend logic
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

		// -- -- Telegram
		if system.System.Config.EnableTelegram {
			telegram.Bot.Init()
			telegram.Bot.SendStartup(strconv.Itoa(system.System.Config.Port))
		}

		go localfolder.ContiniousRead(system.System.Config.OutputFolder)

	} else { // Startup commands = run once and exit
		cmdName := os.Args[1]
		statup_command, exists := command.CMD.Startup.Map[cmdName]
		if !exists {
			system.StartupError()
			prettyprint.P.Cyan.Println("Unknown command:", cmdName)
			log.Println("Unknown command:", cmdName)
			return fmt.Errorf("unknown command")
		}
		statup_command.Run(os.Args[2:]...)
		system.System.Cancel()
	}

	// Nothing should be done here.
	return nil
}
