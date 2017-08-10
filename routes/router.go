package routes

import (
	"github.com/gorilla/mux"
	"net/http"

	"github.com/steffen25/golang.zone/controllers"
	"github.com/steffen25/golang.zone/repositories"
	"github.com/steffen25/golang.zone/database"
	"github.com/steffen25/golang.zone/middlewares"
)

func InitializeRouter(db *database.DB) *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	ur := repositories.UserRepository{DB: db}
	ac := controllers.NewAuthController(&ur)
	uc := controllers.NewUserController(&ur)
	r.HandleFunc("/", middlewares.Logger(uc.HelloWorld)).Methods(http.MethodGet)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/authenticate", middlewares.Logger(ac.Authenticate)).Methods(http.MethodPost)
	api.HandleFunc("/users", middlewares.Logger(uc.GetAll)).Methods(http.MethodGet)
	api.HandleFunc("/users", middlewares.Logger(uc.Create)).Methods(http.MethodPost)
	api.HandleFunc("/users/{id}", middlewares.Logger(uc.GetById)).Methods(http.MethodGet)
	api.HandleFunc("/protected", middlewares.Logger(middlewares.RequireAccessToken(uc.Profile))).Methods(http.MethodGet)

	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/refresh", middlewares.Logger(middlewares.RequireRefreshToken(ac.Refresh))).Methods(http.MethodGet)

	return r
}
