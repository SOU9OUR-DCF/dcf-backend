package config

import (
	"log"
	"sync"
)

var (
	cfg  *Config
	once sync.Once
)

func GetConfig(path string) *Config {
	once.Do(func() {
		var err error
		cfg, err = LoadConfig(path)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
	})
	return cfg
}
