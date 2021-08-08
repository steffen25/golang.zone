package models

import "github.com/go-pg/pg/v10/orm"

func init() {
	orm.RegisterTable((*UserRoles)(nil))
}

type Role struct {
	Id   int
	Name string
}

type UserRoles struct {
	RoleId int
	UserId int
}
