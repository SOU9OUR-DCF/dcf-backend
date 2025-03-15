package config

import (
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	CORS     struct {
		AllowedOrigins []string `yaml:"allowedOrigins"`
	} `yaml:"cors"`
}

type ServerConfig struct {
	Port            string
	Environment     string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

type CookieConfig struct {
	Domain   string
	Path     string
	Secure   bool
	HttpOnly bool
	SameSite string
}
