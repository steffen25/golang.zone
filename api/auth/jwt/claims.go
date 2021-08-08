package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type APIClaims struct {
	UserId          int      `json:"uid"`
	Roles           []string `json:"roles"`
	TokenIdentifier string   `json:"token_identifier"`
	jwt.StandardClaims
}

func (c *APIClaims) ParseClaims(claims jwt.MapClaims) error {
	rl, ok := claims["roles"]
	if !ok {
		return errors.New("could not parse claims roles")
	}

	var roles []string
	if rl != nil {
		for _, v := range rl.([]string) {
			roles = append(roles, v)
		}
	}
	c.Roles = roles

	return nil
}

func (c APIClaims) Authorized(roles ...string) bool {
	for _, has := range c.Roles {
		for _, want := range roles {
			if has == want {
				return true
			}
		}
	}

	return false
}
