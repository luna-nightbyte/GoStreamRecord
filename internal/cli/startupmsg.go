package cli

import (
	"GoStreamRecord/internal/cli/color"
	"GoStreamRecord/internal/db"
	"fmt"
)

func PrintStartup() {

	color.Println("Bblue", `
  ____      ____  _                            ____                        _ 
 / ___| ___/ ___|| |_ _ __ ___  __ _ _ __ ___ |  _ \ ___  ___ ___  _ __ __| |
| |  _ / _ \___ \| __| '__/ _ \/ _' | '_ ' _ \| |_) / _ \/ __/ _ \| '__/ _' |
| |_| | (_) |__) | |_| | |  __/ (_| | | | | | |  _ <  __/ (_| (_) | | | (_| |
 \____|\___/____/ \__|_|  \___|\__,_|_| |_| |_|_| \_\___|\___\___/|_|  \__,_|

	 `)
	color.Println("Bwhite", db.Version)
	color.Println("yellow", "🔹 Written in Go — Fast. Reliable. Efficient.")
	color.Println("yellow", "🔹 Manage streamers, users, and API keys.")
	color.Println("yellow", "🔹 Record what you want, when you want.")
	color.Println("yellow", "🔹 API Ready. Automation Friendly.")
	fmt.Println()
	color.Println("cyan", "📂 Docs: https://github.com/luna-nightbyte/GoStreamRecord")
}
