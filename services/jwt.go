package services

import (
	"time"
	"github.com/steffen25/golang.zone/models"
	"github.com/dgrijalva/jwt-go"
	"log"

	"github.com/steffen25/golang.zone/config"
)

type TokenClaims struct {
	jwt.StandardClaims
	UID int `json:"id"`
}

const (
	tokenDuration = time.Hour * 72
)

func GenerateJWT(u *models.User) (string, error) {
	authClaims := TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenDuration).Unix(),
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