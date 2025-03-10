package prettyprint

import "github.com/fatih/color"

var (
	Cyan      = color.New(color.FgCyan).SprintFunc()
	Green     = color.New(color.FgGreen).SprintFunc()
	Yellow    = color.New(color.FgYellow).SprintFunc()
	BoldRed   = color.New(color.FgRed, color.Bold).SprintFunc()
	BoldWhite = color.New(color.FgWhite, color.Bold).SprintFunc()
	BoldBlue  = color.New(color.FgBlue, color.Bold).SprintFunc()
)
