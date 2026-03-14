package config

import (
	"log"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Env string `koanf:"env"`
	Port int `koanf:"port"`
	MongoURI string `koanf:"mongo_uri"`
}

func Load() *Config {
	k := koanf.New(".")

	err := k.Load(env.Provider("APP_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "APP_"))
	}), nil)

	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	// Set some default values if none is provided
	if cfg.Port == 0 {
		cfg.Port = 8080
	}

	if cfg.Env == "" {
		cfg.Env = "development"
	}

	if cfg.MongoURI == "" {
		cfg.MongoURI = "mongodb://127.0.0.1:27017"
	}

	return &cfg 
}