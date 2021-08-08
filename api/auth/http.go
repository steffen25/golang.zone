package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/qor/roles"
	"github.com/steffen25/golang.zone/api/auth/models"
	"github.com/steffen25/golang.zone/api/auth/rbac"
	"github.com/steffen25/golang.zone/api/middleware"
	"net/http"
)

type AuthRouter struct {
	rbac.Enforcer
	svc AuthService
}

func NewRouter(svc AuthService, rbac rbac.RBAC, public *gin.RouterGroup, private *gin.RouterGroup) {
	ar := AuthRouter{svc: svc}

	loginGroup := public.Group("auth")
	authGroup := private.Group("auth")

	loginGroup.POST("login", ar.authenticate)
	loginGroup.POST("register", ar.register)

	permission := roles.Allow(roles.CRUD, "admin").
		Deny(roles.Create, "guest").
		Allow(roles.Read, "guest")

	authGroup.GET("me", middleware.RBAC(rbac, permission, roles.Read), ar.me)
}

func (ar AuthRouter) authenticate(c *gin.Context) {
	type LoginRequestBody struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var req LoginRequestBody
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "bad request",
		})
		return
	}

	tokens, err := ar.svc.Authenticate(c, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tokens,
	})
}

func (ar AuthRouter) register(c *gin.Context) {
	var newUser models.NewUser
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "unable to decode payload",
		})
		return
	}

	user, err := ar.svc.Register(c, newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": errors.Wrap(err, "unable to create user").Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    user,
	})
}

func (ar AuthRouter) me(c *gin.Context) {
	userId, ok := c.MustGet("auth.user.id").(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "unable to retrieve user",
		})
		return
	}

	//if err := ar.EnforceUser(c, userId); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"code":    http.StatusForbidden,
	//		"message": "insufficient permissions",
	//	})
	//	return
	//}

	user, err := ar.svc.Me(c, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "unable to retrieve user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}
