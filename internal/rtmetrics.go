package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"runtime"
	"time"
)

type RunTimeMetrics struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

type runTimeMetrics struct {
	Alloc,
	TotalAlloc,
	BuckHashSys,
	Frees,
	GCCPUFraction,
	GCSys,
	HeapAlloc,
	HeapIdle,
	HeapInuse,
	HeapObjects,
	HeapReleased,
	HeapSys,
	LastGC,
	Lookups,
	MCacheInuse,
	MCacheSys,
	MSpanInuse,
	MSpanSys,
	Mallocs,
	NextGC,
	NumForcedGC,
	NumGC,
	OtherSys,
	PauseTotalNs,
	StackInuse,
	StackSys,
	Sys,
	RandomValue float64 //gauge
	PollCount int64 //counter
}
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

func NewMonitor(buf io.Writer, key string) func() {
	var rtm runtime.MemStats
	var rm RunTimeMetrics
	rm.Gauge = make(map[string]float64)
	rm.Counter = make(map[string]int64)

	var PollCount int64 = 0
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	check := func(err error) {
		if err != nil {
			log.Println(err)
		}
	}

	return func() {
		runtime.ReadMemStats(&rtm)
		PollCount++
		rm.Gauge["Alloc"] = float64(rtm.Alloc)
		rm.Gauge["TotalAlloc"] = float64(rtm.TotalAlloc)
		rm.Gauge["BuckHashSys"] = float64(rtm.BuckHashSys)
		rm.Gauge["Frees"] = float64(rtm.Frees)
		rm.Gauge["GCCPUFraction"] = rtm.GCCPUFraction
		rm.Gauge["GCSys"] = float64(rtm.GCSys)
		rm.Gauge["HeapAlloc"] = float64(rtm.HeapAlloc)
		rm.Gauge["HeapIdle"] = float64(rtm.HeapIdle)
		rm.Gauge["HeapInuse"] = float64(rtm.HeapInuse)
		rm.Gauge["HeapObjects"] = float64(rtm.HeapObjects)
		rm.Gauge["HeapReleased"] = float64(rtm.HeapReleased)
		rm.Gauge["HeapSys"] = float64(rtm.HeapSys)
		rm.Gauge["LastGC"] = float64(rtm.LastGC)
		rm.Gauge["Lookups"] = float64(rtm.Lookups)
		rm.Gauge["MCacheInuse"] = float64(rtm.MCacheInuse)
		rm.Gauge["MCacheSys"] = float64(rtm.MCacheSys)
		rm.Gauge["MSpanInuse"] = float64(rtm.MSpanInuse)
		rm.Gauge["MSpanSys"] = float64(rtm.MSpanSys)
		rm.Gauge["Mallocs"] = float64(rtm.Mallocs)
		rm.Gauge["NextGC"] = float64(rtm.NextGC)
		rm.Gauge["NumForcedGC"] = float64(rtm.NumForcedGC)
		rm.Gauge["NumGC"] = float64(rtm.NumGC)
		rm.Gauge["OtherSys"] = float64(rtm.OtherSys)
		rm.Gauge["PauseTotalNs"] = float64(rtm.PauseTotalNs)
		rm.Gauge["StackInuse"] = float64(rtm.StackInuse)
		rm.Gauge["StackSys"] = float64(rtm.StackSys)
		rm.Gauge["Sys"] = float64(rtm.Sys)
		rm.Gauge["RandomValue"] = r1.Float64()
		rm.Counter["PollCount"] = PollCount

		//var buf bytes.Buffer
		jsonEncoder := json.NewEncoder(buf)
		if key == "" {
			for id, value := range rm.Gauge {
				check(jsonEncoder.Encode(Metrics{ID: id, MType: "gauge", Value: &value}))
			}
			for id, delta := range rm.Counter {
				check(jsonEncoder.Encode(Metrics{ID: id, MType: "counter", Delta: &delta}))
			}
		} else {
			for id, value := range rm.Gauge {
				check(jsonEncoder.Encode(Metrics{ID: id, MType: "gauge", Value: &value,
					Hash: hash(fmt.Sprintf("%s:gauge:%f", id, value), key)}))
			}
			for id, delta := range rm.Counter {
				check(jsonEncoder.Encode(Metrics{ID: id, MType: "counter", Delta: &delta,
					Hash: hash(fmt.Sprintf("%s:counter:%d", id, delta), key)}))
			}
		}
	}
}

func hash(m, k string) string {
	h := hmac.New(sha256.New, []byte(k))
	h.Write([]byte(m))
	dst := h.Sum(nil)

	log.Printf("%s:%x", m, dst)
	return hex.EncodeToString(dst)
}
