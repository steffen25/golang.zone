package middlewares

import (
	"fmt"
	"net/http"

	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/steffen25/golang.zone/app"
	"github.com/steffen25/golang.zone/controllers"
	"github.com/steffen25/golang.zone/services"
	"os"
)

// TODO: Create error struct that we can use instead of calling controllers?
func RequireAuthentication(a *app.App, next http.HandlerFunc, admin bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		t, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(a.Config.JWT.Secret), nil
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
			jti := claims["jti"].(string)
			tokenHash := claims["tokenHash"].(string)
			val, err := a.Redis.Get(tokenHash+"."+jti).Result()
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
			if !admin {
				next(w, r.WithContext(ctx))
				return
			}
			// Check if the user's token has admin true
			isAdmin := claims["admin"].(bool)
			if !isAdmin {
				controllers.NewAPIError(&controllers.APIError{false, "Admin required", http.StatusForbidden}, w)
				return
			}
			next(w, r.WithContext(ctx))
		}
	}
}

func RequireRefreshToken(a *app.App, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		publicKeyFile, err := os.Open(a.Config.JWT.PublicKeyPath)
		if err != nil {
			panic(err)
		}

		pemfileinfo, _ := publicKeyFile.Stat()
		var size int64 = pemfileinfo.Size()
		pembytes := make([]byte, size)

		buffer := bufio.NewReader(publicKeyFile)
		_, err = buffer.Read(pembytes)

		data, _ := pem.Decode([]byte(pembytes))

		publicKeyFile.Close()

		publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)

		if err != nil {
			panic(err)
		}

		rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

		if !ok {
			panic(err)
		}

		t, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return rsaPub, nil
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
			jti := claims["jti"].(string)
			tokenHash := claims["tokenHash"].(string)
			val, err := a.Redis.Get(tokenHash+"."+jti).Result()
			if err != nil || val == "" {
				controllers.NewAPIError(&controllers.APIError{false, "Invalid token", http.StatusUnauthorized}, w)
				return
			}

			uid := int(claims["id"].(float64))
			ctx := services.ContextWithUserId(r.Context(), uid)
			next(w, r.WithContext(ctx))
		}
	}
}
