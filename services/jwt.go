package services

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"log"

	"github.com/steffen25/golang.zone/models"
	"github.com/steffen25/golang.zone/config"
	"github.com/steffen25/golang.zone/database"
	"fmt"
	"github.com/satori/go.uuid"
	"strconv"
)

type TokenClaims struct {
	jwt.StandardClaims
	UID int `json:"id"`
}

const (
	tokenDuration = time.Hour * 24
)

func GenerateJWT(u *models.User) (string, error) {
	uid := strconv.Itoa(u.ID)
	authClaims := TokenClaims{
		jwt.StandardClaims{
			Id: uid+"."+uuid.NewV4().String(),
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

	redis, _ := database.RedisConnection()
	err = redis.Set(authClaims.Id, u.ID, tokenDuration).Err()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return tokenString, nil
}

func ExtractJti(tokenStr string) (string, error) {
	cfg, err := config.Load("config/app.json")
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// check token signing method etc
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["jti"].(string), nil
	}

	return "", err
}