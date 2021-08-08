package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/steffen25/golang.zone/api/auth/jwt"
	"strings"
)

const (
	BEARER_SCHEMA string = "Bearer "
)

func Authenticate(jwtAuth *jwt.JWTAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("authorization")
		splitToken := strings.Split(token, BEARER_SCHEMA)
		if len(splitToken) != 2 {
			c.Abort()
			c.JSON(401, gin.H{
				"code":    401,
				"message": "malformed token",
			})
			return
		}
		token = strings.TrimSpace(splitToken[1])
		claims, err := jwtAuth.ValidateToken(jwt.AccessTokenType, token)
		if err != nil {
			c.Abort()
			c.JSON(401, gin.H{
				"code":    401,
				"message": err.Error(),
			})
			return
		}

		c.Set("auth.user.id", claims.UserId)
		c.Set("auth.claims", claims)
		c.Next()
	}
}
