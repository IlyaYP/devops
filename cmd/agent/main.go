/*
Задание для трека «Go в DevOps»
Разработайте агент по сбору рантайм-метрик и их последующей отправке на сервер по протоколу HTTP. Разработку нужно вести с использованием шаблона.
Агент должен собирать метрики двух типов:

    gauge, тип float64
    counter, тип int64

В качестве источника метрик используйте пакет runtime.
Нужно собирать следующие метрики:

    Имя метрики: "Alloc", тип: gauge
    Имя метрики: "BuckHashSys", тип: gauge
    Имя метрики: "Frees", тип: gauge
    Имя метрики: "GCCPUFraction", тип: gauge
    Имя метрики: "GCSys", тип: gauge
    Имя метрики: "HeapAlloc", тип: gauge
    Имя метрики: "HeapIdle", тип: gauge
    Имя метрики: "HeapInuse", тип: gauge
    Имя метрики: "HeapObjects", тип: gauge
    Имя метрики: "HeapReleased", тип: gauge
    Имя метрики: "HeapSys", тип: gauge
    Имя метрики: "LastGC", тип: gauge
    Имя метрики: "Lookups", тип: gauge
    Имя метрики: "MCacheInuse", тип: gauge
    Имя метрики: "MCacheSys", тип: gauge
    Имя метрики: "MSpanInuse", тип: gauge
    Имя метрики: "MSpanSys", тип: gauge
    Имя метрики: "Mallocs", тип: gauge
    Имя метрики: "NextGC", тип: gauge
    Имя метрики: "NumForcedGC", тип: gauge
    Имя метрики: "NumGC", тип: gauge
    Имя метрики: "OtherSys", тип: gauge
    Имя метрики: "PauseTotalNs", тип: gauge
    Имя метрики: "StackInuse", тип: gauge
    Имя метрики: "StackSys", тип: gauge
    Имя метрики: "Sys", тип: gauge

К метрикам пакета runtime добавьте другие:

    Имя метрики: "PollCount", тип: counter — счётчик, увеличивающийся на 1 при каждом обновлении метрики из пакета runtime (на каждый pollInterval — см. ниже).
    Имя метрики: "RandomValue", тип: gauge — обновляемое рандомное значение.

По умолчанию приложение должно обновлять метрики из пакета runtime с заданной частотой: pollInterval — 2 секунды.
По умолчанию приложение должно отправлять метрики на сервер с заданной частотой: reportInterval — 10 секунд.
Метрики нужно отправлять по протоколу HTTP, методом POST:

    по умолчанию на адрес: 127.0.0.1, порт: 8080
    в формате: http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
    "application-type": "text/plain"

Агент должен штатно завершаться по сигналам: syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT.
*/
package main

import (
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
	messages := make(chan string, 200)

	go internal.NewMonitor(pollInterval, messages)
	go func() {
		for {
			select {
			case str := <-messages:
				go internal.Send("http://localhost:8080/update" + str)
			default:
				time.Sleep(reportInterval)
			}
		}
	}()
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	log.Println("Shutdown Agent ...")
  os.Exit(0)
}
