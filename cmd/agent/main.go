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

type config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PoolInterval   int    `env:"POLL_INTERVAL"`
}

var cfg config

func init() {
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "Server address")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "Report interval in seconds")
	flag.IntVar(&cfg.PoolInterval, "p", 2, "Poll interval in seconds")
}

func main() {
	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	log.Println("Agent start using args:ADDRESS", cfg.Address, "REPORT_INTERVAL",
		cfg.ReportInterval, "POLL_INTERVAL", cfg.PoolInterval)

	pollInterval := time.Duration(cfg.PoolInterval) * time.Second
	//reportInterval := time.Duration(cfg.ReportInterval) * time.Second
	reportInterval := time.Duration(10) * time.Second
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
}
