package main

import "github.com/steffen25/golang.zone/app"

func main() {
	app := app.New()
	app.Initialize()
	app.Run()
}
