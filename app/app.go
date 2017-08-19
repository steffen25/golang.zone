package app

import (
	"log"
	"net/http"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/steffen25/golang.zone/database"
	"github.com/steffen25/golang.zone/config"
	"github.com/steffen25/golang.zone/routes"
)

type App struct {
	Config 		config.Config
	Database 	*database.MySQLDB
	Redis 	*database.RedisDB
	Router 		*mux.Router
}

func New(cfg config.Config) *App {
	db, err := database.NewMySQLDB(cfg.MySQL)
	if err != nil {
		log.Fatal(err)
	}
	redis, err := database.NewRedisDB(cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}
	router := routes.NewRouter(db)

	return &App{cfg, db, redis, router}
}

func (a *App) Run()  {
	port := a.Config.Port
	addr := fmt.Sprintf(":%v", port)
	fmt.Printf("APP is listening on port: %d\n", port)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) IsProd() bool {
	return a.Config.Env == "prod"
}