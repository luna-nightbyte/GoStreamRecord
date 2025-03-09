package main

import (
	"GoStreamRecord/internal/bot"
	"GoStreamRecord/internal/db"
	"GoStreamRecord/internal/handlers"
	"GoStreamRecord/internal/handlers/cookies"
	"GoStreamRecord/internal/handlers/login"
	"GoStreamRecord/internal/logger"
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"

	"github.com/joho/godotenv"
)

// Embed static HTML files
//
//go:embed internal/app/web/index.html
var IndexHTML string

//go:embed internal/app/web/login.html
var LoginHTML string

func init() {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	boldWhite := color.New(color.FgWhite, color.Bold).SprintFunc()
	boldBlue := color.New(color.FgBlue, color.Bold).SprintFunc()

	fmt.Print(boldBlue(`
  ____      ____  _                            ____                        _ 
 / ___| ___/ ___|| |_ _ __ ___  __ _ _ __ ___ |  _ \ ___  ___ ___  _ __ __| |
| |  _ / _ \___ \| __| '__/ _ \/ _' | '_ ' _ \| |_) / _ \/ __/ _ \| '__/ _' |
| |_| | (_) |__) | |_| | |  __/ (_| | | | | | |  _ <  __/ (_| (_) | | | (_| |
 \____|\___/____/ \__|_|  \___|\__,_|_| |_| |_|_| \_\___|\___\___/|_|  \__,_|

	 `))

	fmt.Println(green("🚀 GoStreamRecorder - ") + boldWhite(db.Version+"\n"))
	fmt.Println(yellow("🔹 Written in Go — Fast. Reliable. Efficient."))
	fmt.Println(yellow("🔹 Manage streamers, users, and API keys."))
	fmt.Println(yellow("🔹 Record what you want, when you want."))
	fmt.Println(yellow("🔹 API Ready. Automation Friendly."))
	fmt.Println()
	fmt.Println(cyan("📂 Docs: https://github.com/luna-nightbyte/GoStreamRecord"))
	handlers.IndexHTML = IndexHTML
	handlers.LoginHTML = LoginHTML

	os.Mkdir("./output", 0755)
	cookies.Session = cookies.New()
	logger.Init(logger.Log_path)
	bot.Init()

}

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] != "reset-pwd" {
			fmt.Println("Usage: ./GoStreamRecord reset-pwd <username> <new-password>")
			fmt.Println("Otherwise run the server without any arguments.")
			return
		}

		if len(os.Args) <= 2 {
			fmt.Println("No username provided.")
			fmt.Println("Usage: ./GoStreamRecord reset-pwd <username> <new-password>")
			fmt.Println("Otherwise run the server without any arguments.")
			return
		}

		username := os.Args[2]
		if len(os.Args) <= 3 {
			fmt.Println("No new password provided.")
			fmt.Println("Usage: ./GoStreamRecord reset-pwd <username> <new-password>")
			fmt.Println("Otherwise run the server without any arguments.")
			return
		}

		newPassword := os.Args[3]
		userFound := false

		for i, u := range db.Config.Users.Users {
			if u.Name == username {
				db.Config.Users.Users[i].Key = string(login.HashedPassword(newPassword))
				userFound = true
				break
			}
		}

		if !userFound {
			log.Println("No matching username found.")
			fmt.Println("No matching username found.")
			return
		}
		db.Config.Update("users", "users.json", &db.Config.Users)
		log.Println("Password updated for", username)
		fmt.Println("Password updated for", username)
		return // Exit after resetting password

	}

	//http.Handle("/", fs)
	handlers.Handle()
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", db.Config.Settings.App.Port),
	}

	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run the server in a separate goroutine
	go func() {
		log.Printf("Server starting on http://127.0.0.1:%d", db.Config.Settings.App.Port)
		fmt.Printf("Server starting on http://127.0.0.1:%d\n", db.Config.Settings.App.Port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for a termination signal
	<-stop
	log.Println("Shutting down server...")
	bot.Bot.StopBot("")
	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited gracefully")
}

var startup_message string = `

`
