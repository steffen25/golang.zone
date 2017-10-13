package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strconv"
	"time"

	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/satori/go.uuid"
	"github.com/steffen25/golang.zone/app"
	"github.com/steffen25/golang.zone/config"
	"github.com/steffen25/golang.zone/models"
	"os"
)

type TokenClaims struct {
	jwt.StandardClaims
	UID   int  `json:"id"`
	Admin bool `json:"admin"`
}

type userCtxKeyType string

const (
	tokenDuration                       = time.Hour
	refreshTokenDuration                = time.Hour * 72
	userCtxKey           userCtxKeyType = "user"
	userIdCtxKey         userCtxKeyType = "userId"
)

type JWTAuthService struct {
	secret     string
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func GenerateJWT(a *app.App, u *models.User) (string, error) {
	uid := strconv.Itoa(u.ID)
	authClaims := TokenClaims{
		jwt.StandardClaims{
			Id:        uid + "." + uuid.NewV4().String(),
			ExpiresAt: time.Now().Add(tokenDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		u.ID,
		u.Admin,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)

	tokenString, err := token.SignedString([]byte(a.Config.JWT.Secret))
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

// TODO: make something like this https://github.com/brainattica/golang-jwt-authentication-api-sample/blob/master/core/authentication/jwt_backend.go
func GenerateRefreshToken(a *app.App, u *models.User) (string, error) {
	privateKeyFile, err := os.Open(a.Config.JWT.PrivateKeyPath)
	if err != nil {
		panic(err)
		return "", err
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
		return "", err
	}
	uid := strconv.Itoa(u.ID)
	authClaims := TokenClaims{
		jwt.StandardClaims{
			Id:        uid + "." + uuid.NewV4().String(),
			ExpiresAt: time.Now().Add(refreshTokenDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		u.ID,
		u.Admin,
	}

	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims = authClaims
	tokenString, err := token.SignedString(privateKeyImported)
	if err != nil {
		panic(err)
		return "", err
	}

	/*uJson, err := json.Marshal(u)
	  if err != nil {
	  	log.Fatal(err)
	  	return "", err
	  }*/
	err = a.Redis.Set(authClaims.Id, u.ID, refreshTokenDuration).Err()
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
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["jti"].(string), nil
	}

	return "", err
}

func ExtractRefreshTokenJti(cfg *config.Config, tokenStr string) (string, error) {
	publicKeyFile, err := os.Open(cfg.JWT.PublicKeyPath)
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

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// check token signing method etc
		return rsaPub, nil
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

			return []byte(cfg.JWT.Secret), nil
		})

	if err != nil || !token.Valid {
		return "", err
	}

	return token.Raw, nil

}

func GetRefreshTokenFromRequest(cfg *config.Config, r *http.Request) (string, error) {
	publicKeyFile, err := os.Open(cfg.JWT.PublicKeyPath)
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
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return rsaPub, nil
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
