package main

import (
	"GoStreamRecord/internal/bot"
	"GoStreamRecord/internal/cli"
	"GoStreamRecord/internal/cli/color"
	"GoStreamRecord/internal/db"
	"GoStreamRecord/internal/handlers"
	"GoStreamRecord/internal/handlers/cookies"
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
)

// Embed static HTML files
//
//go:embed internal/app/web/index.html
var IndexHTML string

//go:embed internal/app/web/login.html
var LoginHTML string

var (
	password_was_reset bool
)

func init() {
	cli.PrintStartup()
	handlers.IndexHTML = IndexHTML
	handlers.LoginHTML = LoginHTML

}

func main() {
	if len(os.Args) < 2 {
		cli.PrintStartup()
		cookies.Session = cookies.New()
		logger.Init(logger.Log_path)
		bot.Init()
		server() // No arguments: run the server.
		return
	}

	cmdName := os.Args[1]
	cmd, exists := cli.Commands[cmdName]
	if !exists {
		fmt.Println()
		color.Print("red", "Unknown command: ")
		color.Println("grey", cmdName)
		// TODO: cli.PrintUsage()
		cli.PrintUsage()
		return
	}

	// Execute the command with the remaining arguments.
	cmd.Execute(os.Args[2:])

}

func server() {
	log.Println("Startup!")
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
