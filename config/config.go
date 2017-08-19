package config

import (
	"os"
	"encoding/json"
	"log"
)

type MySQLConfig struct {
	Username     	string `json:"username"`
	Password 		string `json:"password"`
	DatabaseName    string `json:"database"`
	Encoding        string `json:"encoding"`
}

type RedisConfig struct {
	Host     		string 	`json:"host"`
	Post     		int 	`json:"port"`
}

type Config struct {
	Env string 				`json:"env"`
	MySQL MySQLConfig 		`json:"mysql"`
	Redis RedisConfig 		`json:"redis"`
	Port int 				`json:"port"`
	JWTSecret string 		`json:"jwt_secret"`
}

// New creates a new config by reading a json file that matches the types above
func New(path string) (Config, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	// These are used in different packages so instead of reading the cfg over and over we put them in os.env
	os.Setenv("jwt_secret", cfg.JWTSecret)

	return cfg, nil
}