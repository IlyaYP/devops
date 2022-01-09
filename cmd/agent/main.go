package main

import (
	"bytes"
	"github.com/IlyaYP/devops/internal"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	Address        string `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval int    `env:"REPORT_INTERVAL" envDefault:"10"`
	PoolInterval   int    `env:"POLL_INTERVAL" envDefault:"2"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	pollInterval := time.Duration(cfg.PoolInterval) * time.Second
	reportInterval := time.Duration(cfg.ReportInterval) * time.Second
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
			if err := internal.SendBuf(endPoint, &buf); err != nil {
				log.Fatal(err)
			}
		case <-quit:
			log.Println("Shutdown Agent ...")
			break breakFor
		}
	}
}
