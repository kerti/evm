package config

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/kerti/evm/04-tennis-player/util/logger"
)

var (
	appName = "EVM Machine Gun"
	conf    Config
	once    sync.Once
)

// Config is the configuration struct
type Config struct {
	DB struct {
		Host      string `envconfig:"DB_HOST"`
		Port      int    `envconfig:"DB_PORT"`
		User      string `envconfig:"DB_USER"`
		Pass      string `envconfig:"DB_PASS"`
		Name      string `envconfig:"DB_NAME"`
		ConnLimit int    `envconfig:"DB_CONN_LIMIT"`
	}
	Server struct {
		Port           int           `envconfig:"SERVER_PORT" default:"8080"`
		ShutdownPeriod time.Duration `envconfig:"SERVER_SHUTDOWN_PERIOD" default:"5s"`
	}
}

// Get returns the singleton config instance.
func Get() *Config {
	once.Do(func() {
		err := envconfig.Process("", &conf)
		if err != nil {
			logger.Fatal("Failed to load config: ", err)
		}
		byteConfig, err := json.MarshalIndent(conf, "", "\t")
		if err != nil {
			logger.Fatal("Failed to marshal config: ", err)
		}
		logger.Trace("Config used: %s", byteConfig)
	})
	return &conf
}
