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
	// Repositories
	ur := repositories.NewUserRespository(db)
	pr := repositories.NewPostRepository(db)

	// Controllers
	ac := controllers.NewAuthController(ur)
	uc := controllers.NewUserController(ur, pr)
	pc := controllers.NewPostController(pr)

	r.HandleFunc("/", middlewares.Logger(uc.HelloWorld)).Methods(http.MethodGet)

	api := r.PathPrefix("/api/v1").Subrouter()
	// Users
	api.HandleFunc("/users", middlewares.Logger(uc.GetAll)).Methods(http.MethodGet)
	api.HandleFunc("/users", middlewares.Logger(uc.Create)).Methods(http.MethodPost)
	api.HandleFunc("/users/{id}", middlewares.Logger(uc.GetById)).Methods(http.MethodGet)
	api.HandleFunc("/users/{id}/posts", middlewares.Logger(uc.FindPostsByUser)).Methods(http.MethodGet)
	api.HandleFunc("/protected", middlewares.Logger(middlewares.RequireAuthentication(uc.Profile, false))).Methods(http.MethodGet)

	// Posts
	api.HandleFunc("/posts", middlewares.Logger(pc.GetAll)).Methods(http.MethodGet)
	api.HandleFunc("/posts", middlewares.Logger(middlewares.RequireAuthentication(pc.Create, true))).Methods(http.MethodPost)

	// Authentication
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", middlewares.Logger(ac.Authenticate)).Methods(http.MethodPost)
	auth.HandleFunc("/refresh", middlewares.Logger(middlewares.RequireAuthentication(ac.RefreshToken, false))).Methods(http.MethodGet)
	auth.HandleFunc("/update", middlewares.Logger(middlewares.RequireAuthentication(uc.Update, false))).Methods(http.MethodPut)
	auth.HandleFunc("/logout", middlewares.Logger(middlewares.RequireAuthentication(ac.Logout, false))).Methods(http.MethodGet)
	auth.HandleFunc("/logout/all", middlewares.Logger(middlewares.RequireAuthentication(ac.LogoutAll, false))).Methods(http.MethodGet)

	return r
}
