package services

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"log"
	"strconv"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/satori/go.uuid"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/steffen25/golang.zone/models"
	"github.com/steffen25/golang.zone/config"
	"github.com/steffen25/golang.zone/app"
)

type TokenClaims struct {
	jwt.StandardClaims
	UID int `json:"id"`
	Admin bool `json:"admin"`
}

const (
	tokenDuration = time.Hour * 24
	userCtxKey userCtxKeyType = "user"
	userIdCtxKey userCtxKeyType = "userId"
)

type userCtxKeyType string

func GenerateJWT(a *app.App, u *models.User) (string, error) {
	uid := strconv.Itoa(u.ID)
	authClaims := TokenClaims{
		jwt.StandardClaims{
			Id: uid+"."+uuid.NewV4().String(),
			ExpiresAt: time.Now().Add(tokenDuration).Unix(),
			IssuedAt: time.Now().Unix(),
		},
		u.ID,
		u.Admin,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)

	tokenString, err := token.SignedString([]byte(a.Config.JWTSecret))
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	/*uJson, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
		return "", err
	}*/
	err = a.Redis.Set(authClaims.Id, u.ID, tokenDuration).Err()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return tokenString, nil
}

func ExtractJti(cfg *config.Config, tokenStr string) (string, error) {
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

func GetTokenFromRequest(cfg *config.Config, r *http.Request) (string, error) {
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(cfg.JWTSecret), nil
		})

	if err != nil || !token.Valid {
		return "", err
	}

	return token.Raw, nil

}

// TODO: https://www.calhoun.io/pitfalls-of-context-values-and-how-to-avoid-or-mitigate-them/
func ContextWithUserId(ctx context.Context, uID int) context.Context {
	return context.WithValue(ctx, userIdCtxKey, uID)
}

func UserIdFromContext(ctx context.Context) (int, error) {
	uID, ok := ctx.Value(userIdCtxKey).(int)
	if !ok {
		log.Println("Context missing userID")
		return -1, errors.New("[SERVICE]: Context missing userID")
	}

	return uID, nil
}

func ContextWithUser(ctx context.Context, u *models.User) context.Context {
	return context.WithValue(ctx, userCtxKey, u)
}

func UserFromContext(ctx context.Context) (*models.User, error) {
	u, ok := ctx.Value(userCtxKey).(*models.User)
	if !ok {
		log.Println("Context missing user")
		return nil, errors.New("[SERVICE]: Context missing userID")
	}

	return u, nil
}