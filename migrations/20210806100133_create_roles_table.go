package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
			CREATE TABLE roles (
				id SERIAL PRIMARY KEY,
				name varchar(255) NOT NULL,
				description varchar(255) NOT NULL
			)
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("DROP TABLE roles")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210806100133_create_roles_table", up, down, opts)
}
