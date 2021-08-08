package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
			CREATE TABLE users (
				id SERIAL PRIMARY KEY,
				name varchar(255) NOT NULL,
				email varchar(255) NOT NULL UNIQUE,
				password varchar(255) NOT NULL,
				created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
				deleted_at TIMESTAMPTZ DEFAULT NULL
			)
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("DROP TABLE users")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20201213024256_create_users_table", up, down, opts)
}
