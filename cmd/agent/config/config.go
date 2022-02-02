package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

//type config struct {
//	Address        string `env:"ADDRESS" envDefault:"localhost:8080"`
//	ReportInterval int    `env:"REPORT_INTERVAL" envDefault:"10"`
//	PoolInterval   int    `env:"POLL_INTERVAL" envDefault:"2"`
//}

//type config struct {
//	Address        string `env:"ADDRESS"`
//	ReportInterval int    `env:"REPORT_INTERVAL"`
//	PoolInterval   int    `env:"POLL_INTERVAL"`
//}

type Config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PoolInterval   time.Duration `env:"POLL_INTERVAL"`
	Key            string        `env:"KEY"`
	EndPoint       string
}

func LoadConfig() (*Config, error) {
	var cfg Config
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "Server address")
	flag.DurationVar(&cfg.ReportInterval, "r", time.Duration(10)*time.Second, "Report interval in seconds")
	flag.DurationVar(&cfg.PoolInterval, "p", time.Duration(2)*time.Second, "Poll interval in seconds")
	flag.StringVar(&cfg.Key, "k", "", "Key")
	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		log.Println(err)
		return nil, err
	}

	//cfg.EndPoint = "http://" + cfg.Address + "/update/"
	cfg.EndPoint = "http://" + cfg.Address + "/updates/"

	return &cfg, nil
}
