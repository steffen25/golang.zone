package main

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/kelseyhightower/envconfig"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
	"github.com/steffen25/golang.zone/api"
	"log"
	"os"
)

func main() {
	const migrationDir = "migrations"
	var cfg api.Config
	err := envconfig.Process("golangzone", &cfg)
	if err != nil {
		log.Fatalf("error processing1 config %v", err)
	}

	dbAddr := fmt.Sprintf("%s:%s", cfg.DB.Host, cfg.DB.Port)
	db := pg.Connect(&pg.Options{
		Addr:     dbAddr,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Database: cfg.DB.Database,
	})
	ctx := context.Background()

	if err := db.Ping(ctx); err != nil {
		panic(err)
	}
	err = migrations.Run(db, migrationDir, os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
