package api

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-redis/redis/v8"
	"github.com/kelseyhightower/envconfig"
	"github.com/steffen25/golang.zone/api/auth/jwt"
	"github.com/steffen25/golang.zone/api/handlers"
	"log"
)

type Config struct {
	DB struct {
		Host     string `envconfig:"DB_HOST"`
		Port     string `envconfig:"DB_PORT"`
		Database string `envconfig:"DB_DATABASE"`
		User     string `envconfig:"DB_USERNAME"`
		Password string `envconfig:"DB_PASSWORD"`
	}
	Redis struct {
		Host     string `envconfig:"REDIS_HOST"`
		Port     string `envconfig:"REDIS_PORT"`
		Database int    `envconfig:"REDIS_DATABASE"`
		Password string `envconfig:"REDIS_PASSWORD"`
	}
	Auth struct {
		AccessTokenAlgorithm  string `envconfig:"JWT_ACCESS_ALGORITHM"`
		RefreshTokenAlgorithm string `envconfig:"JWT_REFRESH_ALGORITHM"`
		JWTSecret             string `envconfig:"JWT_SECRET_KEY"`
		JWTPublicKey          string `envconfig:"JWT_PUBLIC_KEY"`
		JWTPrivateKey         string `envconfig:"JWT_PRIVATE_KEY"`
	}
}

func New() {
	var cfg Config
	err := envconfig.Process("golangzone", &cfg)
	if err != nil {
		log.Fatalf("error processing config %v", err)
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
		log.Fatal(err)
	}

	redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",                 // no password set
		DB:       cfg.Redis.Database, // use default DB
	})

	rdbCtx := context.Background()
	if status := rdb.Ping(rdbCtx); status.Err() != nil {
		log.Fatal(status.Err())
	}

	aCfg := jwt.JWTConfig(cfg.Auth)

	jwtAuth, err := jwt.New(aCfg, rdb, db)
	if err != nil {
		log.Fatalf("error constructing auth %v", err)
	}

	r := handlers.API(db, rdb, jwtAuth)
	//pprof.Register(r)

	r.Run(":8080")
}
