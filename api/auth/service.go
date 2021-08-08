package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/steffen25/golang.zone/api/auth/jwt"
	"github.com/steffen25/golang.zone/api/auth/models"
)

type AuthService interface {
	Authenticate(c *gin.Context, email string, password string) (jwt.TokenPair, error)
	Register(c *gin.Context, newUser models.NewUser) (models.User, error)
	Me(c *gin.Context, userId int) (models.User, error)
}

type TokenGenerator interface {
	GenerateTokens(accessClaims, refreshClaims jwt.APIClaims) (jwt.TokenPair, error)
}

type AuthStore interface {
	FindByEmail(email string) (models.User, error)
	FindById(userId int) (models.User, error)
	CreateUser(user *models.User) error
}

type Auth struct {
	db             *pg.DB
	Store          AuthStore
	tokenGenerator TokenGenerator
	rbac           RBAC
}

type RBAC interface {
	EnforceUser(*gin.Context, int) error
}

func New(db *pg.DB, authStore AuthStore, tokenGenerator TokenGenerator, rbac RBAC) Auth {
	return Auth{db: db, Store: authStore, tokenGenerator: tokenGenerator, rbac: rbac}
}
