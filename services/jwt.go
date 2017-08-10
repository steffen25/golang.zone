package services

import (
	"time"
	"github.com/steffen25/golang.zone/models"
	"github.com/dgrijalva/jwt-go"
	"log"

	"github.com/steffen25/golang.zone/config"
	"github.com/satori/go.uuid"
	"context"
)

type TokenClaims struct {
	jwt.StandardClaims
	UID int `json:"id"`
}

const (
	ACCESS_TOKEN_DURATION = time.Minute * 30
	REFRESH_TOKEN_DURATION = time.Hour * 72
)

func GenerateJWT(u *models.User) (string, error) {
	authClaims := TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ACCESS_TOKEN_DURATION).Unix(),
			IssuedAt: time.Now().Unix(),
		},
		u.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)

	// TODO: Find a better way to pass the config from the App(redundant)
	cfg, err := config.Load("config/app.json")
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(u *models.User) (string, error) {
	authClaims := TokenClaims{
		jwt.StandardClaims{
			Id: uuid.NewV4().String(),
			ExpiresAt: time.Now().Add(REFRESH_TOKEN_DURATION).Unix(),
			IssuedAt: time.Now().Unix(),
		},
		u.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)

	// TODO: Find a better way to pass the config from the App(redundant)
	cfg, err := config.Load("config/app.json")
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return tokenString, nil
}

func NewContextWithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, "user", user)
}

func UserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value("user").(*models.User)
	return user, ok
}