package models

import (
	"crypto/rand"
	"fmt"
	gojwt "github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	uuid "github.com/satori/go.uuid"
	"github.com/steffen25/golang.zone/api/auth/jwt"
	"strconv"
	"time"
)

type User struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	Password  string      `json:"-"`
	Roles     []Role      `pg:"many2many:user_roles" json:"roles"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	DeletedAt pg.NullTime `pg:"deleted_at,soft_delete" json:"deleted_at,omitempty"`
}

// NewUser represents the struct used to create a new user on golang.zone
type NewUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// BeforeInsert hook executed before database insert operation.
func (u *User) BeforeInsert(db orm.DB) error {
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
		u.UpdatedAt = now
	}
	return nil
	//return a.Validate()
}

// BeforeUpdate hook executed before database update operation.
func (u *User) BeforeUpdate(db orm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
	//return a.Validate()
}

func (u User) RoleNames() []string {
	var names []string

	for _, r := range u.Roles {
		names = append(names, r.Name)
	}

	return names
}

func (u *User) Claims() (access, refresh jwt.APIClaims, error error) {
	uId := strconv.Itoa(u.ID)
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return jwt.APIClaims{}, jwt.APIClaims{}, err
	}
	tokenId := fmt.Sprintf("%x", key)

	var userRoles []string

	for _, role := range u.Roles {
		userRoles = append(userRoles, role.Name)
	}

	nowUnix := time.Now().Unix()

	return jwt.APIClaims{
			StandardClaims: gojwt.StandardClaims{
				Audience:  "",
				ExpiresAt: time.Now().Add(jwt.AccessTokenDuration).Unix(),
				Id:        uId + "." + uuid.NewV4().String(),
				IssuedAt:  nowUnix,
				Issuer:    "golang.zone",
				NotBefore: 0,
				Subject:   "",
			},
			UserId:          u.ID,
			Roles:           userRoles, // store roles in jwt? Middleware will always check the database.
			TokenIdentifier: tokenId,
		}, jwt.APIClaims{
			StandardClaims: gojwt.StandardClaims{
				Audience:  "",
				ExpiresAt: time.Now().Add(jwt.RefreshTokenDuration).Unix(),
				Id:        uId + "." + uuid.NewV4().String(),
				IssuedAt:  nowUnix,
				Issuer:    "golang.zone",
				NotBefore: 0,
				Subject:   "",
			},
			UserId:          u.ID,
			Roles:           userRoles,
			TokenIdentifier: tokenId,
		}, nil
}
