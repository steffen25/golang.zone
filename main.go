package main

import (
	"log"
	
	"github.com/steffen25/golang.zone/app"
	"github.com/steffen25/golang.zone/config"
	"github.com/steffen25/golang.zone/routes"
)

func main() {
	cfg, err := config.New("config/app.json")
	if err != nil {
		log.Fatal(err)
	}
	app := app.New(cfg)
	router := routes.NewRouter(app)
	app.Run(router)
}