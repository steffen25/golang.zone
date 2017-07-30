package middlewares

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/steffen25/golang.zone/config"
	"github.com/steffen25/golang.zone/controllers"
)

func RequireJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg, err := config.Load("config/app.json")
		if err != nil {
			log.Fatal(err)
		}

		t, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(cfg.JWTSecret), nil
			})

		if err != nil || t == nil {
			w.WriteHeader(http.StatusUnauthorized)
			err := controllers.NewAPIError(false, "Invalid JWT", http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err)
			return
		}

		next(w, r)
	}
}