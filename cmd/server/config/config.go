package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

// move defaults to flag Parse
//type Config struct {
//	Address       string `env:"ADDRESS" envDefault:"localhost:8080"`
//	StoreInterval int    `env:"STORE_INTERVAL" envDefault:"300"`
//	StoreFile     string `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
//	Restore       bool   `env:"RESTORE" envDefault:"true"`
//}

//type Config struct {
//	Address       string `env:"ADDRESS"`
//	StoreInterval int    `env:"STORE_INTERVAL"`
//	StoreFile     string `env:"STORE_FILE"`
//	Restore       bool   `env:"RESTORE"`
//}

const (
	defaultConfigEndpoint = "localhost:8080"
)

type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	Key           string        `env:"KEY"`
	DbDsn         string        `env:"DATABASE_DSN"`
}

// Validate performs a basic validation.
func (c Config) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("%s field: empty", "Address")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		Address: defaultConfigEndpoint,
	}
}

func LoadConfig() (*Config, error) {
	var cfg Config
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "Server address")
	flag.DurationVar(&cfg.StoreInterval, "i", time.Duration(300)*time.Second, "Store interval in seconds")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "Store file")
	flag.BoolVar(&cfg.Restore, "r", true, "Restore data from file when start")
	flag.StringVar(&cfg.Key, "k", "", "Key")
	flag.StringVar(&cfg.DbDsn, "d", "", "DATABASE_DSN")
	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		log.Println(err)
		return nil, err
	}

	return &cfg, nil
}
