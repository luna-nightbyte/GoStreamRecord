package startup

import (
	"GoStreamRecord/internal/db"
	"GoStreamRecord/internal/prettyprint"
	"fmt"
)

func PrintUsage() {
	fmt.Println(prettyprint.Cyan("Usage:"))
	for _, cmd := range Commands {
		fmt.Println(prettyprint.Cyan(" - " + cmd.Usage))
	}
	fmt.Println(prettyprint.Cyan("Otherwise run the server without any arguments."))
}

func PrintStartup() {

	fmt.Print(prettyprint.BoldBlue(`
  ____      ____  _                            ____                        _ 
 / ___| ___/ ___|| |_ _ __ ___  __ _ _ __ ___ |  _ \ ___  ___ ___  _ __ __| |
| |  _ / _ \___ \| __| '__/ _ \/ _' | '_ ' _ \| |_) / _ \/ __/ _ \| '__/ _' |
| |_| | (_) |__) | |_| | |  __/ (_| | | | | | |  _ <  __/ (_| (_) | | | (_| |
 \____|\___/____/ \__|_|  \___|\__,_|_| |_| |_|_| \_\___|\___\___/|_|  \__,_|

	 `))

	fmt.Println(prettyprint.Green("🚀 GoStreamRecorder - ") + prettyprint.BoldWhite(db.Version+"\n"))
	fmt.Println(prettyprint.Yellow("🔹 Written in Go — Fast. Reliable. Efficient."))
	fmt.Println(prettyprint.Yellow("🔹 Manage streamers, users, and API keys."))
	fmt.Println(prettyprint.Yellow("🔹 Record what you want, when you want."))
	fmt.Println(prettyprint.Yellow("🔹 API Ready. Automation Friendly."))
	fmt.Println()
	fmt.Println(prettyprint.Cyan("📂 Docs: https://github.com/luna-nightbyte/GoStreamRecord"))
}
