package prettyprint

import (
	"fmt"

	"github.com/fatih/color"
)

// var (
// 	Cyan      = color.New(color.FgCyan).SprintFunc()
// 	Green     = color.New(color.FgGreen).SprintFunc()
// 	Yellow    = color.New(color.FgYellow).SprintFunc()
// 	BoldRed   = color.New(color.FgRed, color.Bold).SprintFunc()
// 	BoldWhite = color.New(color.FgWhite, color.Bold).SprintFunc()
// 	BoldGrey  = color.New(color.FgHiBlack, color.Bold).SprintFunc()
// 	BoldBlue  = color.New(color.FgBlue, color.Bold).SprintFunc()
// 	BoldGreen = color.New(color.FgGreen, color.Bold).SprintFunc()
// )

type colors struct {
	Cyan      iColor
	Green     iColor
	Yellow    iColor
	BoldRed   iColor
	BoldWhite iColor
	BoldGrey  iColor
	BoldBlue  iColor
	BoldGreen iColor
}

type iColor struct {
	c func(a ...interface{}) string
}

var P colors

func init() {
	P.Cyan.c = color.New(color.FgCyan).SprintFunc()
	P.Green.c = color.New(color.FgGreen).SprintFunc()
	P.Yellow.c = color.New(color.FgYellow).SprintFunc()
	P.BoldRed.c = color.New(color.FgRed, color.Bold).SprintFunc()
	P.BoldWhite.c = color.New(color.FgWhite, color.Bold).SprintFunc()
	P.BoldGrey.c = color.New(color.FgHiBlack, color.Bold).SprintFunc()
	P.BoldBlue.c = color.New(color.FgBlue, color.Bold).SprintFunc()
	P.BoldGreen.c = color.New(color.FgGreen, color.Bold).SprintFunc()

}
func (xc *iColor) Println(a ...any) {
	fmt.Println(xc.c(a...))
}

func (xc *iColor) Print(a ...any) {
	fmt.Print(xc.c(a...))
}

func (xc *iColor) Color(a ...any) string {
	return xc.c(a)
}
