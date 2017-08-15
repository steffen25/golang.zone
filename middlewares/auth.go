package middlewares

import (
	"net/http"
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/steffen25/golang.zone/config"
	"github.com/steffen25/golang.zone/controllers"
	"github.com/steffen25/golang.zone/database"
	"github.com/steffen25/golang.zone/services"
)

// TODO: Create error struct that we can use instead of calling controllers?
func RequireAuthentication(next http.HandlerFunc) http.HandlerFunc {
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

		if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
			redis, _ := database.RedisConnection()
			jti := claims["jti"].(string)
			val, err := redis.Get(jti).Result()
			if err != nil || val == "" {
				controllers.NewAPIError(&controllers.APIError{false, "Invalid token", http.StatusUnauthorized}, w)
				return
			}
			// TODO: Put the user into the context instead of the user id? Right now we only need to reference to the id of the user that is logged in
			// maybe put the json representation of the user inside redis and use the 'val' here
			/*user := &models.User{}
			err = json.Unmarshal([]byte(val), &user)
			if err != nil {
				controllers.NewAPIError(&controllers.APIError{false, "Something went wrong", http.StatusInternalServerError}, w)
				return
			}
			ctx := services.ContextWithUser(r.Context(), user)*/
			uid := int(claims["id"].(float64))
			/*db, err := database.NewDB(cfg.Database)
			if err != nil {
				controllers.NewAPIError(&controllers.APIError{false, "Something went wrong", http.StatusInternalServerError}, w)
				return
			}
			userRepo := repositories.NewUserRespository(db)
			user, err := userRepo.FindById(uid)
			if err != nil {
				controllers.NewAPIError(&controllers.APIError{false, "Something went wrong", http.StatusInternalServerError}, w)
				return
			}
			ctx := services.ContextWithUser(r.Context(), user)*/
			ctx := services.ContextWithUserId(r.Context(), uid)
			next(w, r.WithContext(ctx))
		}
	}
}