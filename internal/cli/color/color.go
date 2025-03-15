package color

import (
	"fmt"

	c "github.com/fatih/color"
)

var Colors = map[string]func(a ...interface{}) string{
	"cyan":    c.New(c.FgCyan).SprintFunc(),
	"green":   c.New(c.FgGreen).SprintFunc(),
	"yellow":  c.New(c.FgYellow).SprintFunc(),
	"red":     c.New(c.FgRed).SprintFunc(),
	"white":   c.New(c.FgWhite).SprintFunc(),
	"blue":    c.New(c.FgBlue).SprintFunc(),
	"Bcyan":   c.New(c.FgCyan).SprintFunc(),
	"Bgreen":  c.New(c.FgGreen).SprintFunc(),
	"grey":    c.New(c.FgBlack).SprintFunc(),
	"Byellow": c.New(c.FgYellow, c.Bold).SprintFunc(),
	"Bred":    c.New(c.FgRed, c.Bold).SprintFunc(),
	"Bwhite":  c.New(c.FgWhite, c.Bold).SprintFunc(),
	"Bblue":   c.New(c.FgBlue, c.Bold).SprintFunc(),
	"Bgrey":   c.New(c.FgBlack, c.Bold).SprintFunc(),
}

// fmt.Print
//
// name = name (Normal letters)
// name = Bname (Bold letters)
func Print(name, msg string) {
	fmt.Print(Colors[name](msg))
}

// fmt.Println
//
// name = name (Normal letters)
// name = Bname (Bold letters)
func Println(name, msg string) {
	fmt.Println(Colors[name](msg))
}
