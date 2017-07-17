package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/steffen25/golang.zone/middlewares"
	"github.com/steffen25/golang.zone/controllers"
)

type Router struct {
	*mux.Router
}


func InitializeRouter() *Router {
	r := mux.NewRouter()
	r.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/users/{id}", middlewares.Logger((&controllers.UserController{}).Get))

	return &Router{r}
}
