package main

import (
	"bytes"
	"flag"
	"github.com/IlyaYP/devops/internal"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
	"os/signal"
	"syscall"
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
type config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PoolInterval   time.Duration `env:"POLL_INTERVAL"`
}

var cfg config

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "Server address")
	flag.DurationVar(&cfg.ReportInterval, "r", time.Duration(10)*time.Second, "Report interval in seconds")
	flag.DurationVar(&cfg.PoolInterval, "p", time.Duration(2)*time.Second, "Poll interval in seconds")
	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		log.Println(err)
		return err
	}
	log.Println("Agent start using args:ADDRESS", cfg.Address, "REPORT_INTERVAL",
		cfg.ReportInterval, "POLL_INTERVAL", cfg.PoolInterval)
	pollInterval := cfg.PoolInterval
	reportInterval := cfg.ReportInterval

	endPoint := "http://" + cfg.Address + "/update/"
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	var buf bytes.Buffer

	getMetrics := internal.NewMonitor(&buf)
	poll := time.Tick(pollInterval)
	report := time.Tick(reportInterval)
breakFor:
	for {
		select {
		case <-poll:
			getMetrics()
		case <-report:
			if err := internal.SendBufRetry(endPoint, &buf); err != nil {
				log.Println(err)
				log.Println("Ok, let's try again later")
			}
		case <-quit:
			log.Println("Shutdown Agent ...")
			break breakFor
		}
	}
	return nil
}
