package config

import (
	"log"
	"os"
	"time"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-deafult:"local"`
	HTTPServer `yaml:"http_server"`
	Database `yaml:"database"`
}

type HTTPServer struct {
	Host 		string `yaml:"host" env-default:"localhost"`
	Port 		string `yaml:"port" env-default:"8080"`
	Timeout 	time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Database struct {
	Host     string `yaml:"host" env-default:"localhost"`
	User     string `yaml:"user" env-default:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name" env-default:"postgres"`
	Port     string `yaml:"port" env-default:"5432"`
	SSLMode  string `yaml:"ssl_mode" env-default:"disable"`
}

func MustLoad() Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
	    log.Fatalf("cannot read config: %s", err)
	}

	return cfg
}


