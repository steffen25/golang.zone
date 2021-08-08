package jwt

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg/v10"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"time"
)

type JWTAuth struct {
	cfg                JWTConfig
	db                 *pg.DB
	rdb                *redis.Client
	accessTokenAlgo    jwt.SigningMethod
	refreshTokenAlgo   jwt.SigningMethod
	accessTokenParser  *jwt.Parser
	refreshTokenParser *jwt.Parser
}

type JWTConfig struct {
	AccessTokenAlgorithm  string
	RefreshTokenAlgorithm string
	JWTSecret             string
	JWTPublicKey          string
	JWTPrivateKey         string
}

type AccessToken struct {
	jwt.StandardClaims
}

type RefreshToken struct {
	jwt.StandardClaims
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
}

type TokenType string

const (
	AccessTokenDuration            = time.Hour
	RefreshTokenDuration           = time.Hour * 24 * 30
	AccessTokenType      TokenType = "access"
	RefreshTokenType     TokenType = "refresh"
)

func New(cfg JWTConfig, rdb *redis.Client, db *pg.DB) (*JWTAuth, error) {
	accessMethod := jwt.GetSigningMethod(cfg.AccessTokenAlgorithm)
	if accessMethod == nil {
		return nil, errors.Errorf("unknown access token algorithm %v", cfg.AccessTokenAlgorithm)
	}

	refreshMethod := jwt.GetSigningMethod(cfg.RefreshTokenAlgorithm)
	if refreshMethod == nil {
		return nil, errors.Errorf("unknown refresh algorithm %v", cfg.RefreshTokenAlgorithm)
	}

	accessTokenParser := jwt.Parser{
		ValidMethods: []string{cfg.AccessTokenAlgorithm},
	}

	refreshTokenParser := jwt.Parser{
		ValidMethods: []string{cfg.RefreshTokenAlgorithm},
	}

	jwtAuth := JWTAuth{
		cfg:                cfg,
		db:                 db,
		rdb:                rdb,
		accessTokenAlgo:    accessMethod,
		refreshTokenAlgo:   refreshMethod,
		accessTokenParser:  &accessTokenParser,
		refreshTokenParser: &refreshTokenParser,
	}

	return &jwtAuth, nil
}

func (jwtAuth *JWTAuth) GenerateTokens(accessClaims, refreshClaims APIClaims) (TokenPair, error) {
	accessToken := jwt.NewWithClaims(jwtAuth.accessTokenAlgo, accessClaims)
	tokenStr, err := accessToken.SignedString([]byte(jwtAuth.cfg.JWTSecret))
	if err != nil {
		return TokenPair{}, errors.Wrap(err, "signing access token")
	}

	// move somewhere else
	data, _ := pem.Decode([]byte(jwtAuth.cfg.JWTPrivateKey))
	privKey, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		return TokenPair{}, errors.Wrap(err, "parsing jwt private key")
	}

	refreshToken := jwt.NewWithClaims(jwtAuth.refreshTokenAlgo, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString(privKey)
	if err != nil {
		return TokenPair{}, errors.Wrap(err, "signing refresh token")
	}

	return TokenPair{
		AccessToken:  tokenStr,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    3600,
		TokenType:    "Bearer",
	}, nil
}

func (jwtAuth *JWTAuth) ValidateToken(tokenType TokenType, token string) (APIClaims, error) {
	var claims APIClaims
	var parser *jwt.Parser
	var keyFunc jwt.Keyfunc

	switch tokenType {
	case AccessTokenType:
		parser = jwtAuth.accessTokenParser
		keyFunc = func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtAuth.cfg.JWTSecret), nil
		}
	case RefreshTokenType:
		parser = jwtAuth.refreshTokenParser
		keyFunc = func(token *jwt.Token) (interface{}, error) {
			data, _ := pem.Decode([]byte(jwtAuth.cfg.JWTPublicKey))
			pubKey, err := x509.ParsePKCS1PublicKey(data.Bytes)
			if err != nil {
				return nil, err
			}
			return pubKey, nil
		}
	default:
		panic("unknown token type")
	}

	authToken, err := parser.ParseWithClaims(token, &claims, keyFunc)
	if err != nil {
		return APIClaims{}, errors.Wrap(err, "parsing token")
	}

	if !authToken.Valid {
		return APIClaims{}, errors.New("invalid token")
	}

	return claims, nil
}
