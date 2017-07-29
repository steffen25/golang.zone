package app

import (
	"log"
	"net/http"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/steffen25/golang.zone/database"
	"github.com/steffen25/golang.zone/config"
	"github.com/gorilla/mux"
	"github.com/steffen25/golang.zone/routes"
)

type App struct {
	Config 		config.Config
	Database 	*database.DB
	Router 		*mux.Router
}

func New() *App {
	return &App{}
	// Call the Initialize here instead of in main? make it not exportable?
}

func (a *App) Initialize()  {
	var err error
	cfg, err := config.Load("config/app.json")
	if err != nil {
		log.Fatal(err)
	}
	a.Config = cfg
	db, err := database.NewDB(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	a.Database = db
	a.Router = routes.InitializeRouter(db)
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