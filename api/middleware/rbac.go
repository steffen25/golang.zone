package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/qor/roles"
	"github.com/steffen25/golang.zone/api/auth/jwt"
	"github.com/steffen25/golang.zone/api/auth/rbac"
	"net/http"
)

func Authorize(rbac rbac.RBAC, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.MustGet("auth.claims").(jwt.APIClaims)
		if !ok {
			c.Abort()
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "you are not authorized",
			})
			return
		}

		// get userid from claims and pass downstream to RBAC to query database for roles.
		err := rbac.CheckRoles(claims.UserId, roles...)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "you are not authorized for that action",
			})
			return
		}
		c.Next()
	}
}

func RBAC(rbac rbac.RBAC, roles *roles.Permission, action roles.PermissionMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.MustGet("auth.claims").(jwt.APIClaims)
		if !ok {
			c.Abort()
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "you are not authorized",
			})
			return
		}

		// get userid from claims and pass downstream to RBAC to query database for roles.
		err := rbac.CheckPermission(claims.UserId, roles, action)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "you are not authorized for that action",
			})
			return
		}
		c.Next()
	}
}
