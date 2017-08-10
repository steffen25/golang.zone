package middlewares

import (
	"net/http"
	"fmt"
	"log"
	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/steffen25/golang.zone/config"
	"github.com/steffen25/golang.zone/controllers"
)

// TODO: Create error struct that we can use instead of calling controllers?
func RequireAccessToken(next http.HandlerFunc) http.HandlerFunc {
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


		if err != nil {
			if err == request.ErrNoTokenInRequest {
				controllers.NewAPIError(&controllers.APIError{false, "Missing token", http.StatusUnauthorized}, w)
				return
			}

			controllers.NewAPIError(&controllers.APIError{false, "Invalid token", http.StatusUnauthorized}, w)
			return
		}

		// TODO: Check if token is blacklisted in redis

		if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
			// Only allow access token to be used
			if _, ok := claims["id"]; !ok {
				controllers.NewAPIError(&controllers.APIError{false, "Invalid token", http.StatusUnauthorized}, w)
				return
			}
			ctx := context.WithValue(r.Context(), "userId", claims["id"])
			next(w, r.WithContext(ctx))
		}
	}
}

func RequireRefreshToken(next http.HandlerFunc) http.HandlerFunc {
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


		if err != nil {
			if err == request.ErrNoTokenInRequest {
				controllers.NewAPIError(&controllers.APIError{false, "Missing token", http.StatusUnauthorized}, w)
				return
			}

			controllers.NewAPIError(&controllers.APIError{false, "Invalid token", http.StatusUnauthorized}, w)
			return
		}

		// TODO: Check if token is blacklisted in redis

		if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
			// Only allow access token to be used
			if _, ok := claims["jti"]; !ok {
				controllers.NewAPIError(&controllers.APIError{false, "Invalid token", http.StatusUnauthorized}, w)
				return
			}
			ctx := context.WithValue(r.Context(), "userId", claims["id"])
			next(w, r.WithContext(ctx))
		}
	}
}