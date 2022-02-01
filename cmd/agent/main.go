package main

import (
	"bytes"
	"github.com/IlyaYP/devops/cmd/agent/config"
	"github.com/IlyaYP/devops/internal"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Agent start using args:ADDRESS", cfg.Address, "REPORT_INTERVAL",
		cfg.ReportInterval, "POLL_INTERVAL", cfg.PoolInterval, "KEY", cfg.Key)

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	var buf bytes.Buffer

	Metrics := internal.NewRTM(&buf, cfg.Key)
	poll := time.Tick(cfg.PoolInterval)
	report := time.Tick(cfg.ReportInterval)
breakFor:
	for {
		select {
		case <-poll:
			Metrics.Collect()
		case <-report:
			if err := internal.SendBufRetry(cfg.EndPoint, Metrics.GetJSON()); err != nil {
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
