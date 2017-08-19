package main

import (
	"github.com/steffen25/golang.zone/app"
	"github.com/steffen25/golang.zone/config"
	"log"
)

func main() {
	cfg, err := config.New("config/app.json")
	if err != nil {
		log.Fatal(err)
	}
	app := app.New(cfg)
	app.Run()
}
