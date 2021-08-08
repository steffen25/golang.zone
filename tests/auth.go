package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/steffen25/golang.zone/api/auth/jwt"
	"github.com/steffen25/golang.zone/api/auth/models"
)

type AuthMock struct {
	FindByEmailFn    func(email string) (models.User, error)
	GenerateTokensFn func(accessClaims, refreshClaims jwt.APIClaims) (jwt.TokenPair, error)
	EnforceUserFn    func(*gin.Context, int) error
}

func (a *AuthMock) EnforceUser(context *gin.Context, i int) error {
	panic("implement me")
}

func (a *AuthMock) FindById(userId int) (models.User, error) {
	panic("implement me")
}

func (a *AuthMock) CreateUser(user *models.User) error {
	panic("implement me")
}

func (a *AuthMock) GenerateTokens(accessClaims, refreshClaims jwt.APIClaims) (jwt.TokenPair, error) {
	return a.GenerateTokensFn(accessClaims, refreshClaims)
}

func (a *AuthMock) FindByEmail(email string) (models.User, error) {
	return a.FindByEmailFn(email)
}
