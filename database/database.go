package database

import (
	"fmt"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/steffen25/golang.zone/config"
)

type DB struct {
	*sql.DB
}

func NewDB(dbCfg config.MysqlConfig) (*DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@/%s", dbCfg.Username, dbCfg.Password, dbCfg.DatabaseName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}