package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
			CREATE TABLE user_roles (
				role_id INT NOT NULL,
				CONSTRAINT fk_role
      				FOREIGN KEY(role_id) 
	  				REFERENCES roles(id),
				user_id INT NOT NULL,
				CONSTRAINT fk_user
      				FOREIGN KEY(user_id) 
	  				REFERENCES users(id),
				PRIMARY KEY(role_id, user_id)
			)
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("DROP TABLE user_roles")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210806100539_create_user_roles_table", up, down, opts)
}
