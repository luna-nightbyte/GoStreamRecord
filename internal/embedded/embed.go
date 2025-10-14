package embedded

import "embed"

// import (
// 	"embed"
// 	_ "embed"
// )

// //go:embed login/dist/index.html
var LoginHTML string

// //go:embed app/dist/index.html
// var VueFrontend string

//go:embed app/dist/*
var VueDistFiles embed.FS
