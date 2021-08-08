package rbac

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/qor/roles"
	"github.com/steffen25/golang.zone/api/auth/models"
	"log"
)

type UserRole uint

const (
	IsAdmin UserRole = 1 << iota
	IsEditor
)

type RBAC interface {
	CheckRoles(userId int, roles ...string) error
	CheckPermission(userId int, roles *roles.Permission, action roles.PermissionMode) error
}

type Enforcer struct{}

type rbac struct {
	Enforcer
	db *pg.DB
}

func New(db *pg.DB) rbac {
	return rbac{db: db}
}

func (rbac rbac) CheckRoles(userId int, roles ...string) error {
	u := models.User{ID: userId}

	err := rbac.db.Model(&u).
		Column("id").
		Where("id = ?id").
		Relation("Roles").
		First()

	if err != nil {
		log.Println("could not find user", err)
		return err
	}

	log.Println("[RBAC] roles:", u.Roles)

	// check roles
	for _, userRoles := range u.Roles {
		for _, want := range roles {
			if userRoles.Name == want {
				return nil
			}
		}
	}

	return errors.New("you do not have the correct permissions")
}

func (rbac rbac) CheckPermission(userId int, roles *roles.Permission, action roles.PermissionMode) error {
	u := models.User{ID: userId}

	// todo: should we inject a db or a authstore? The same query is used within the store.
	err := rbac.db.Model(&u).
		Column("id").
		Where("id = ?id").
		Relation("Roles").
		First()

	if err != nil {
		return err
	}

	var roleNames []interface{}

	for _, r := range u.Roles {
		roleNames = append(roleNames, r.Name)
	}

	ok := roles.HasPermission(action, roleNames...)
	if !ok {
		return errors.New("insufficient permissions")
	}

	return nil
}

// EnforceUser checks whether the id equals the id of the authenticated user.
// The authenticated user id is set by the auth middleware.
func (rbac Enforcer) EnforceUser(c *gin.Context, id int) error {
	// todo: if admin return nil
	uid, ok := c.MustGet("auth.user.id").(int)
	if !ok {
		return errors.New("forbidden")
	}

	if uid == id {
		return nil
	}

	return errors.New("forbidden")
}
