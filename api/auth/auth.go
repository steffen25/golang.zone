package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/steffen25/golang.zone/api/auth/jwt"
	"github.com/steffen25/golang.zone/api/auth/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func (a Auth) Authenticate(c *gin.Context, email, password string) (jwt.TokenPair, error) {
	// check if user exists and check password.
	user, err := a.Store.FindByEmail(email)
	if err != nil {
		log.Println(err)
		return jwt.TokenPair{}, errors.Wrap(err, "finding user")
	}

	// compare password hash in db with password.
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return jwt.TokenPair{}, errors.New("authentication failed")
	}

	ac, rc, err := user.Claims()
	if err != nil {
		return jwt.TokenPair{}, errors.Wrap(err, "generate user jwt claims")
	}

	return a.tokenGenerator.GenerateTokens(ac, rc)
}

func (a Auth) Register(c *gin.Context, newUser models.NewUser) (models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, errors.Wrap(err, "generating password hash")
	}

	user := models.User{
		Name:     newUser.Name,
		Email:    newUser.Email,
		Password: string(hash),
	}

	err = a.Store.CreateUser(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (a Auth) Me(c *gin.Context, userId int) (models.User, error) {
	//if err := a.rbac.EnforceUser(c, userId); err != nil {
	//	return models.User{}, err
	//}
	return a.Store.FindById(userId)
}
