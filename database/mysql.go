package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/steffen25/golang.zone/config"
)

type MySQLDB struct {
	*sql.DB
}

func NewMySQLDB(dbCfg config.MySQLConfig) (*MySQLDB, error) {
	//DSN := fmt.Sprintf("%s:%s@unix(/tmp/mysql.sock)/%s?parseTime=true", dbCfg.Username, dbCfg.Password, dbCfg.DatabaseName)
	dataSourceName := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=true", dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DatabaseName, dbCfg.Encoding)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &MySQLDB{db}, nil
}
