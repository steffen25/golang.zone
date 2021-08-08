package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/go-redis/redis/v8"
	"github.com/steffen25/golang.zone/api/auth"
	"github.com/steffen25/golang.zone/api/auth/jwt"
	"github.com/steffen25/golang.zone/api/auth/rbac"
	"github.com/steffen25/golang.zone/api/database"
	"github.com/steffen25/golang.zone/api/middleware"
	"github.com/steffen25/golang.zone/api/user"
)

func API(db *pg.DB, rdb *redis.Client, jwtAuth *jwt.JWTAuth) *gin.Engine {
	router := gin.Default()

	v1Public := router.Group("api/v1")
	v1Auth := router.Group("api/v1")

	v1Auth.Use(middleware.Authenticate(jwtAuth))

	authStore := database.NewAuthStore(db)

	rbacHelper := rbac.New(db)

	authSvc := auth.New(db, authStore, jwtAuth, rbacHelper)
	auth.NewRouter(authSvc, rbacHelper, v1Public, v1Auth)

	userSvc := user.NewService(db)
	user.NewRouter(userSvc, rbacHelper, v1Auth)

	return router
}
