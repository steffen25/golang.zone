package config

import (
	"os"
	"encoding/json"
	"log"
)

type MysqlConfig struct {
	Username     	string `json:"db_username"`
	Password 		string `json:"db_password"`
	DatabaseName    string `json:"db_database"`
}

type Config struct {
	Env string 				`json:"env"`
	Database MysqlConfig 	`json:"database"`
	Port int 				`json:"port"`
}



// LoadConfig will read config/app.json and parse its content into the Config type
func Load(path string) (cfg Config, err error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg, nil
}