package user

import (
	"github.com/gin-gonic/gin"
	"github.com/steffen25/golang.zone/api/auth/rbac"
)

type UserRouter struct {
	svc UserService
}

func NewRouter(svc UserService, rbac rbac.RBAC, r *gin.RouterGroup) {
	ur := UserRouter{svc}

	group := r.Group("users")

	group.POST("", ur.create)
}

func (ur UserRouter) create(c *gin.Context) {
	ur.svc.Create()
}
