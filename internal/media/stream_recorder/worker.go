package stream_recorder

import (
	"context"
	"fmt"
	"log"
	"remoteCtrl/internal/media/stream_recorder/recorder"
	"strings"
	"sync"
)

// Bot encapsulates the recording bot’s state.
type Controller struct {
	mux        sync.Mutex
	status     []recorder.Recorder
	isFirstRun bool
	logger     *log.Logger
	// ctx is used to signal shutdown.
	ctx    context.Context
	cancel context.CancelFunc
}

var Streamer Controller = NewBot()

// NewBot creates a new Bot, sets up its cancellation context.
func NewBot() Controller {
	ctx, cancel := context.WithCancel(context.Background())
	b := Controller{
		ctx:        ctx,
		cancel:     cancel,
		status:     []recorder.Recorder{},
		isFirstRun: true,
	}
	return b
}

func (b *Controller) Execute(command string, name string) {

	if len(command) == 0 {
		log.Println("No command provided..")
		return
	}
	switch strings.ToLower(command) {
	case "start":
		fmt.Println("Starting bot for", name)
		// If the bot was previously stopped, reinitialize the context.
		if b.ctx.Err() != nil {
			b.ctx, b.cancel = context.WithCancel(context.Background())
		}

		for _, s := range b.status {
			if name == s.Website.Username && s.Cmd != nil {
				fmt.Println("Bot already running for", name)
				log.Println("Bot already running..")
				return
			}
		}
		fmt.Println("Starting bot")
		log.Println("Starting bot")
		go b.RecordLoop(name)
	case "stop":
		is_running := len(b.status) != 0

		if !is_running && len(b.status) == 0 {
			log.Println("[bot] Stopped recording for")
			break
		}

		log.Println("Stopping bot")
		var wg sync.WaitGroup
		b.mux.Lock()
		// Iterate over a copy of the status slice to avoid closure capture issues.
		for _, s := range b.status {
			// Stop only the specified process (or all if name is empty).
			if name == "" || s.Website.Username == name {

				b.stopProcessIfRunning(s)
				sName := s.Website.Username
				wg.Add(1)
				go func(n string) {
					defer wg.Done()
					b.StopProcess(n)
				}(sName)
			} else {
				log.Println("Not stopping..")
			}
		}

		b.mux.Unlock()
		wg.Wait()

		b.checkProcesses()
	case "restart":
		log.Println("Restarting bot")
		recorders := []string{}
		// Before restarting, reinitialize the context so RecordLoop doesn't exit immediately.
		b.ctx, b.cancel = context.WithCancel(context.Background())

		if name != "" {
			// Stop a single process.
			process := getProcess(name, b)

			b.Execute("stop", process.Website.Username)
			recorders = append(recorders, name)

		} else {
			var wg sync.WaitGroup
			// Stop all running recorders.
			// Create a copy of b.status to avoid data races when stopping processes.
			b.mux.Lock()
			statusCopy := make([]recorder.Recorder, len(b.status))
			copy(statusCopy, b.status)
			b.mux.Unlock()
			for _, s := range statusCopy {
				b.mux.Lock()

				// Mark that the process is being restarted.
				// (Assuming b.status is the source of truth; you might also update the copy)
				for i, rec := range b.status {
					if rec.Website.Username == s.Website.Username {
						b.status[i].IsRestarting = true
						b.stopProcessIfRunning(b.status[i])
						break
					}
				}
				b.mux.Unlock()
				wg.Add(1)
				recorders = append(recorders, s.Website.Username)
				go func(n string) {
					b.Execute("stop", n)
					log.Println("Stopped", n)
					wg.Done()
				}(s.Website.Username)
			}
			wg.Wait()
		}

		// Start all recorders that were stopped.
		for _, recName := range recorders {
			go b.RecordLoop(recName)
		}
	default:
		log.Println("Nothing to do..")
	}
}
