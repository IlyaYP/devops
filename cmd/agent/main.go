package main

import (
	"bytes"
	"context"
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	var buf bytes.Buffer

	Metrics := internal.NewRTM(&buf, cfg.Key)
	Metrics.PoolInterval = cfg.PoolInterval
	go Metrics.Run(ctx)
	//poll := time.Tick(cfg.PoolInterval)
	report := time.Tick(cfg.ReportInterval)
breakFor:
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutdown Agent ...")
			break breakFor
		//case <-poll:
		//	Metrics.Collect()
		case <-report:
			//log.Println("before", buf.Len())
			Metrics.GetJSON()
			//log.Println("after", buf.Len())
			if err := internal.SendBufRetry(cfg.EndPoint, &buf); err != nil {
				log.Println(err)
				log.Println("Ok, let's try again later")
			}
			//log.Println("after send", buf.Len())
		}
	}
	return nil
}
