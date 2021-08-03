package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/steffen25/golang.zone/app"
	"github.com/steffen25/golang.zone/controllers"
	"github.com/steffen25/golang.zone/middlewares"
	"github.com/steffen25/golang.zone/repositories"
	"github.com/steffen25/golang.zone/services"
)

func NewRouter(a *app.App) *mux.Router {
	r := mux.NewRouter()
	// Repositories
	ur := repositories.NewUserRespository(a.Database)
	pr := repositories.NewPostRepository(a.Database)

	// Services
	jwtAuth := services.NewJWTAuthService(&a.Config.JWT, a.Redis)

	// Controllers
	ac := controllers.NewAuthController(a, ur, jwtAuth)
	uc := controllers.NewUserController(a, ur, pr)
	pc := controllers.NewPostController(a, pr, ur)
	uploadController := controllers.NewUploadController()

	r.HandleFunc("/", middlewares.Logger(uc.HelloWorld)).Methods(http.MethodGet)

	api := r.PathPrefix("/api/v1").Subrouter()

	// Uploads
	api.HandleFunc("/images/upload", middlewares.Logger(middlewares.RequireAuthentication(a, uploadController.UploadImage, true))).Methods(http.MethodPost)

	// Users
	api.HandleFunc("/users", middlewares.Logger(uc.GetAll)).Methods(http.MethodGet)
	api.HandleFunc("/users/{id}", middlewares.Logger(uc.GetById)).Methods(http.MethodGet)
	api.HandleFunc("/users/{id}/posts", middlewares.Logger(uc.FindPostsByUser)).Methods(http.MethodGet)
	api.HandleFunc("/protected", middlewares.Logger(middlewares.RequireAuthentication(a, uc.Profile, false))).Methods(http.MethodGet)

	// Posts
	api.HandleFunc("/posts", middlewares.Logger(pc.GetAll)).Methods(http.MethodGet)
	api.HandleFunc("/posts/{id:[0-9]+}", middlewares.Logger(pc.GetById)).Methods(http.MethodGet)
	api.HandleFunc("/posts/{slug}", middlewares.Logger(pc.GetBySlug)).Methods(http.MethodGet)
	api.HandleFunc("/posts", middlewares.Logger(middlewares.RequireAuthentication(a, pc.Create, true))).Methods(http.MethodPost)
	api.HandleFunc("/posts/{id}", middlewares.Logger(middlewares.RequireAuthentication(a, pc.Update, true))).Methods(http.MethodPut)

	// Authentication
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", middlewares.Logger(ac.Register)).Methods(http.MethodPost)
	auth.HandleFunc("/login", middlewares.Logger(ac.Authenticate)).Methods(http.MethodPost)
	auth.HandleFunc("/me", middlewares.Logger(middlewares.RequireAuthentication(a, ac.Me, false))).Methods(http.MethodGet)
	auth.HandleFunc("/refresh", middlewares.Logger(middlewares.RequireRefreshToken(a, ac.RefreshTokens))).Methods(http.MethodGet)
	auth.HandleFunc("/update", middlewares.Logger(middlewares.RequireAuthentication(a, uc.Update, false))).Methods(http.MethodPut)
	auth.HandleFunc("/logout", middlewares.Logger(middlewares.RequireAuthentication(a, ac.Logout, false))).Methods(http.MethodGet)
	auth.HandleFunc("/logout/all", middlewares.Logger(middlewares.RequireAuthentication(a, ac.LogoutAll, false))).Methods(http.MethodGet)

	return r
}
