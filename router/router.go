package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/steffen25/golang.zone/controllers"
	"github.com/steffen25/golang.zone/repositories"
	"github.com/steffen25/golang.zone/database"
	"github.com/steffen25/golang.zone/middlewares"
)

type Router struct {
	*mux.Router
}

func InitializeRouter(db *database.DB) *Router {
	r := mux.NewRouter()
	r.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	ur := repositories.UserRepository{DB: db}
	uc := controllers.NewUserController(&ur)
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/hello", middlewares.Logger(uc.HelloWorld)).Methods(http.MethodGet)
	api.HandleFunc("/users", middlewares.Logger(uc.GetAll)).Methods(http.MethodGet)
	api.HandleFunc("/users", middlewares.Logger(uc.Create)).Methods(http.MethodPost)
	api.HandleFunc("/users/{id}", middlewares.Logger(uc.GetById)).Methods(http.MethodGet)

	return &Router{r}
}
