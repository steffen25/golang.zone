package services

import (
	"bufio"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/satori/go.uuid"
	"github.com/steffen25/golang.zone/config"
	"github.com/steffen25/golang.zone/database"
	"github.com/steffen25/golang.zone/models"
	"github.com/steffen25/golang.zone/util"
)

type TokenClaims struct {
	jwt.StandardClaims
	UID   int  `json:"id"`
	Admin bool `json:"admin"`
}

type KAuthTokenClaims struct {
	jwt.StandardClaims
	UID   int  `json:"id"`
	Admin bool `json:"admin"`
	TokenHash string `json:"tokenHash"`
}

type AccessToken struct {
	AccessToken string `json:"accessToken"`
}

type RefreshToken struct {
	RefreshToken string `json:"refreshToken"`
}

type Tokens struct {
	AccessToken  string  `json:"accessToken"`
	RefreshToken string  `json:"refreshToken"`
	ExpiresIn    float64 `json:"expiresIn"`
	TokenType    string  `json:"tokenType"`
}

type userCtxKeyType string

const (
	TokenDuration                       = time.Hour
	RefreshTokenDuration                = time.Hour * 72
	TokenType                           = "Bearer"
	userCtxKey           userCtxKeyType = "user"
	userIdCtxKey         userCtxKeyType = "userId"
)

type JWTAuthService interface {
	GenerateTokens(u *models.User) (*Tokens, error)
}

type jwtAuthService struct {
	secret     string
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	Redis      *database.RedisDB
}

func NewJWTAuthService(jwtCfg *config.JWTConfig, redis *database.RedisDB) JWTAuthService {
	return &jwtAuthService{
		jwtCfg.Secret,
		getPrivateKey(jwtCfg),
		getPublicKey(jwtCfg),
		redis,
	}
}

func (jwtService *jwtAuthService) GenerateTokens(u *models.User) (*Tokens, error) {
	uid := strconv.Itoa(u.ID)
	now := time.Now()
	tokenHash := util.GetMD5Hash(now.String()+uid)
	authClaims := KAuthTokenClaims{
		jwt.StandardClaims{
			Id:        uid + "." + uuid.NewV4().String(),
			ExpiresAt: now.Add(TokenDuration).Unix(),
			IssuedAt:  now.Unix(),
		},
		u.ID,
		u.Admin,
		tokenHash,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)

	accessTokenString, err := accessToken.SignedString([]byte(jwtService.secret))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = jwtService.Redis.Set(tokenHash + "." + authClaims.Id, u.ID, TokenDuration).Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodRS512)
	authClaims.Id = uid + "." + uuid.NewV4().String()
	refreshToken.Claims = authClaims
	refreshTokenString, err := refreshToken.SignedString(jwtService.privateKey)
	if err != nil {
		panic(err)
		return nil, err
	}

	err = jwtService.Redis.Set(tokenHash + "." + authClaims.Id, u.ID, RefreshTokenDuration).Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	tokens := &Tokens{
		accessTokenString,
		refreshTokenString,
		3600,
		TokenType,
	}

	return tokens, nil
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

func ExtractTokenHash(cfg *config.Config, tokenStr string) (string, error) {
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
		return claims["tokenHash"].(string), nil
	}

	return "", err
}

func ExtractRefreshTokenHash(cfg *config.Config, tokenStr string) (string, error) {
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
		return claims["tokenHash"].(string), nil
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
		return nil, errors.New("[SERVICE]: Context missing user")
	}

	return u, nil
}

func getPrivateKey(jwtConfig *config.JWTConfig) *rsa.PrivateKey {
	privateKeyFile, err := os.Open(jwtConfig.PrivateKeyPath)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	return privateKeyImported
}

func getPublicKey(jwtConfig *config.JWTConfig) *rsa.PublicKey {
	publicKeyFile, err := os.Open(jwtConfig.PublicKeyPath)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	return rsaPub
}
