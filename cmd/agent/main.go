package main

import (
	"bytes"
	"github.com/IlyaYP/devops/internal"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	pollInterval := time.Duration(2) * time.Second
	reportInterval := time.Duration(10) * time.Second
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
			if err := internal.SendBuf("http://localhost:8080/update/", &buf); err != nil {
				log.Fatal(err)
			}
		case <-quit:
			log.Println("Shutdown Agent ...")
			break breakFor
		}
	}
}
