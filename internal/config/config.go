package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Token string `yml:"token"`
}

func MustLoad() *Config {
	configPath := os.Getenv("HO4UHA_BOT_CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yml"
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config %s don't exists", configPath)
	}
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatal("Can't read config file", err)
	}
	return &cfg
}
