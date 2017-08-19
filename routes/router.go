package routes

import (
	"github.com/gorilla/mux"
	"net/http"

	"github.com/steffen25/golang.zone/controllers"
	"github.com/steffen25/golang.zone/repositories"
	"github.com/steffen25/golang.zone/middlewares"
	"github.com/steffen25/golang.zone/app"
)

func NewRouter(a *app.App) *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	// Repositories
	ur := repositories.NewUserRespository(a.Database)
	pr := repositories.NewPostRepository(a.Database)

	// Controllers
	ac := controllers.NewAuthController(a, ur)
	uc := controllers.NewUserController(a, ur, pr)
	pc := controllers.NewPostController(a, pr)

	r.HandleFunc("/", middlewares.Logger(uc.HelloWorld)).Methods(http.MethodGet)

	api := r.PathPrefix("/api/v1").Subrouter()
	// Users
	api.HandleFunc("/users", middlewares.Logger(uc.GetAll)).Methods(http.MethodGet)
	api.HandleFunc("/users", middlewares.Logger(uc.Create)).Methods(http.MethodPost)
	api.HandleFunc("/users/{id}", middlewares.Logger(uc.GetById)).Methods(http.MethodGet)
	api.HandleFunc("/users/{id}/posts", middlewares.Logger(uc.FindPostsByUser)).Methods(http.MethodGet)
	api.HandleFunc("/protected", middlewares.Logger(middlewares.RequireAuthentication(a, uc.Profile, false))).Methods(http.MethodGet)

	// Posts
	api.HandleFunc("/posts", middlewares.Logger(pc.GetAll)).Methods(http.MethodGet)
	api.HandleFunc("/posts", middlewares.Logger(middlewares.RequireAuthentication(a, pc.Create, true))).Methods(http.MethodPost)
	api.HandleFunc("/posts/{id}", middlewares.Logger(middlewares.RequireAuthentication(a, pc.Update, true))).Methods(http.MethodPut)

	// Authentication
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", middlewares.Logger(ac.Authenticate)).Methods(http.MethodPost)
	auth.HandleFunc("/refresh", middlewares.Logger(middlewares.RequireAuthentication(a, ac.RefreshToken, false))).Methods(http.MethodGet)
	auth.HandleFunc("/update", middlewares.Logger(middlewares.RequireAuthentication(a, uc.Update, false))).Methods(http.MethodPut)
	auth.HandleFunc("/logout", middlewares.Logger(middlewares.RequireAuthentication(a, ac.Logout, false))).Methods(http.MethodGet)
	auth.HandleFunc("/logout/all", middlewares.Logger(middlewares.RequireAuthentication(a, ac.LogoutAll, false))).Methods(http.MethodGet)

	return r
}
