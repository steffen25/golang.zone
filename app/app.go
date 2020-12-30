package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/steffen25/golang.zone/config"
	"github.com/steffen25/golang.zone/database"
)

type App struct {
	Config   config.Config
	Database *database.MySQLDB
	Redis    *database.RedisDB
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

	return &App{cfg, db, redis}
}

func (a *App) Run(r *mux.Router) {
	var corsOptions []handlers.CORSOption
	if !a.IsProd() {
		creds := handlers.AllowCredentials()
		corsOptions = append(corsOptions, creds)
	}
	corsOptions = append(corsOptions, handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "X-Requested-With"}))
	corsOptions = append(corsOptions, handlers.AllowedOrigins([]string{"*"}))
	corsOptions = append(corsOptions, handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}))
	port := a.Config.Port
	addr := fmt.Sprintf(":%v", port)
	fmt.Printf("APP is listening on port: %d\n", port)
	log.Fatal(http.ListenAndServe(addr, handlers.CORS(corsOptions...)(r)))
}

func (a *App) IsProd() bool {
	return a.Config.Env == "prod"
}
