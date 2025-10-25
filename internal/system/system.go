package system

import (
	"context"
	"os"
	"remoteCtrl/internal/db"
)

type Core struct {
	IsOnline        bool
	WaitForNetwork  bool
	Config          db.Config
	Context         context.Context
	Cancel          context.CancelFunc
	triggerShutdown chan os.Signal
}

var (
	System        Core
	onlineCheckIP string = "8.8.8.8"
	enableDebug   string = "false"
	DEBUG                = enableDebug == "true"
)
