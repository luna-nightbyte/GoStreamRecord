package main

import "embed"

// import (
// 	"embed"
// 	_ "embed"
// )

// //go:embed login/dist/index.html
// var LoginHTML string

//go:embed vue/login/dist/*
var VueLoginFiles embed.FS

// //go:embed app/dist/index.html
// var VueFrontend string

//go:embed vue/app/dist/*
var VueDistFiles embed.FS
