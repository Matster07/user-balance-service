package configs

import (
	"github.com/caarlos0/env"
	"github.com/matster07/user-balance-service/internal/pkg/logging"
	"sync"
)

type Config struct {
	Port             int    `env:"PORT" envDefault:"3000"`
	ApiVersion       string `env:"API_VERSION" envDefault:"v1"`
	DatabaseHost     string `env:"DATABASE_HOST" envDefault:"127.0.0.1"`
	DatabasePort     int    `env:"DATABASE_PORT" envDefault:"5432"`
	DatabaseTable    string `env:"DATABASE_TABLE" envDefault:"postgres"`
	DatabaseUsername string `env:"DATABASE_USERNAME" envDefault:"postgres"`
	DatabasePassword string `env:"DATABASE_PASSWORD" envDefault:"1926Semul!"`
	DatabaseSchema   string `env:"DATABASE_SCHEMA" envDefault:"user-balance-service"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}
		if err := env.Parse(instance); err != nil {
			logger.Fatal(err)
		}
	})
	return instance
}
