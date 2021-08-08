package user

import "github.com/go-pg/pg/v10"

type UserService interface {
	Create() error
}

type User struct {
	db *pg.DB
}

func NewService(db *pg.DB) *User {
	return &User{db: db}
}
