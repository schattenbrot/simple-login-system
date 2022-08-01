package config

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	Port int
	Env  string
	DB   struct {
		DSN string
	}
	JWT []byte
}

type Application struct {
	Version   string
	StartTime time.Time
	Config    Config
	Logger    *log.Logger
	Validator *validator.Validate
}
