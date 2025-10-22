package prettyprint

import (
	"fmt"

	"github.com/fatih/color"
)

// The colors struct holds all the color configurations as iColor structs.
// We've expanded this to include standard colors, high-intensity variants (Light),
// style attributes (Underline, Faint), and combined styles (Success, BackgroundBlue).
type colors struct {
	// --- Original Colors ---
	Cyan      iColor
	Green     iColor
	Yellow    iColor
	BoldRed   iColor
	BoldWhite iColor
	BoldGrey  iColor // Uses FgHiBlack
	BoldBlue  iColor
	BoldGreen iColor

	// --- New Standard Colors ---
	Red     iColor
	Magenta iColor
	White   iColor

	// --- New High-Intensity Colors (Light) ---
	LightRed     iColor
	LightGreen   iColor
	LightYellow  iColor
	LightBlue    iColor
	LightMagenta iColor
	LightCyan    iColor
	LightWhite   iColor

	// --- New Style Combinations ---
	UnderlineWhite iColor
	FaintWhite     iColor
	Success        iColor // Bold Black text on Green background
	BackgroundBlue iColor // White text on Blue background
}

// iColor is a wrapper around the color printing function.
type iColor struct {
	c func(a ...interface{}) string
}

// P is the global instance through which all color functions are accessed (e.g., P.Cyan.Println("Hello")).
var P colors

// init initializes all the color functions for the global P variable.
func init() {
	// --- Original Colors ---
	P.Cyan.c = color.New(color.FgCyan).SprintFunc()
	P.Green.c = color.New(color.FgGreen).SprintFunc()
	P.Yellow.c = color.New(color.FgYellow).SprintFunc()
	P.BoldRed.c = color.New(color.FgRed, color.Bold).SprintFunc()
	P.BoldWhite.c = color.New(color.FgWhite, color.Bold).SprintFunc()
	P.BoldGrey.c = color.New(color.FgHiBlack, color.Bold).SprintFunc() // HiBlack usually renders as bright grey
	P.BoldBlue.c = color.New(color.FgBlue, color.Bold).SprintFunc()
	P.BoldGreen.c = color.New(color.FgGreen, color.Bold).SprintFunc()

	// --- New Standard Colors ---
	P.Red.c = color.New(color.FgRed).SprintFunc()
	P.Magenta.c = color.New(color.FgMagenta).SprintFunc()
	P.White.c = color.New(color.FgWhite).SprintFunc()

	// --- New High-Intensity Colors (Light) ---
	P.LightRed.c = color.New(color.FgHiRed).SprintFunc()
	P.LightGreen.c = color.New(color.FgHiGreen).SprintFunc()
	P.LightYellow.c = color.New(color.FgHiYellow).SprintFunc()
	P.LightBlue.c = color.New(color.FgHiBlue).SprintFunc()
	P.LightMagenta.c = color.New(color.FgHiMagenta).SprintFunc()
	P.LightCyan.c = color.New(color.FgHiCyan).SprintFunc()
	P.LightWhite.c = color.New(color.FgHiWhite).SprintFunc()

	// --- New Style Combinations ---
	P.UnderlineWhite.c = color.New(color.FgWhite, color.Underline).SprintFunc()
	P.FaintWhite.c = color.New(color.FgWhite, color.Faint).SprintFunc()
	P.Success.c = color.New(color.FgBlack, color.BgGreen, color.Bold).SprintFunc()
	P.BackgroundBlue.c = color.New(color.FgWhite, color.BgBlue).SprintFunc()
}

// Println prints the given arguments followed by a newline, using the color of the receiver.
func (xc *iColor) Println(a ...any) {
	fmt.Println(xc.c(a...))
}

// Print prints the given arguments using the color of the receiver.
func (xc *iColor) Print(a ...any) {
	fmt.Print(xc.c(a...))
}

// Color returns the colored string without printing it, suitable for embedded formatting.
func (xc *iColor) Color(a ...any) string {
	return xc.c(a...)
}
