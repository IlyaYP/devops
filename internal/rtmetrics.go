package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
		for id, value := range rm.Gauge {
			check(jsonEncoder.Encode(Metrics{ID: id, MType: "gauge", Value: &value}))
		}
		for id, delta := range rm.Counter {
			check(jsonEncoder.Encode(Metrics{ID: id, MType: "counter", Delta: &delta}))
		}
		/*
			check(jsonEncoder.Encode(Metrics{ID: "Alloc", MType: "gauge", Value: &rm.Alloc}))
			check(jsonEncoder.Encode(Metrics{ID: "TotalAlloc", MType: "gauge", Value: &rm.TotalAlloc}))
			check(jsonEncoder.Encode(Metrics{ID: "BuckHashSys", MType: "gauge", Value: &rm.BuckHashSys}))
			check(jsonEncoder.Encode(Metrics{ID: "Frees", MType: "gauge", Value: &rm.Frees}))
			check(jsonEncoder.Encode(Metrics{ID: "GCCPUFraction", MType: "gauge", Value: &rm.GCCPUFraction}))
			check(jsonEncoder.Encode(Metrics{ID: "GCSys", MType: "gauge", Value: &rm.GCSys}))
			check(jsonEncoder.Encode(Metrics{ID: "HeapAlloc", MType: "gauge", Value: &rm.HeapAlloc}))
			check(jsonEncoder.Encode(Metrics{ID: "HeapIdle", MType: "gauge", Value: &rm.HeapIdle}))
			check(jsonEncoder.Encode(Metrics{ID: "HeapInuse", MType: "gauge", Value: &rm.HeapInuse}))
			check(jsonEncoder.Encode(Metrics{ID: "HeapObjects", MType: "gauge", Value: &rm.HeapObjects}))
			check(jsonEncoder.Encode(Metrics{ID: "HeapReleased", MType: "gauge", Value: &rm.HeapReleased}))
			check(jsonEncoder.Encode(Metrics{ID: "HeapSys", MType: "gauge", Value: &rm.HeapSys}))
			check(jsonEncoder.Encode(Metrics{ID: "LastGC", MType: "gauge", Value: &rm.LastGC}))
			check(jsonEncoder.Encode(Metrics{ID: "Lookups", MType: "gauge", Value: &rm.Lookups}))
			check(jsonEncoder.Encode(Metrics{ID: "MCacheInuse", MType: "gauge", Value: &rm.MCacheInuse}))
			check(jsonEncoder.Encode(Metrics{ID: "MCacheSys", MType: "gauge", Value: &rm.MCacheSys}))
			check(jsonEncoder.Encode(Metrics{ID: "MSpanInuse", MType: "gauge", Value: &rm.MSpanInuse}))
			check(jsonEncoder.Encode(Metrics{ID: "MSpanSys", MType: "gauge", Value: &rm.MSpanSys}))
			check(jsonEncoder.Encode(Metrics{ID: "Mallocs", MType: "gauge", Value: &rm.Mallocs}))
			check(jsonEncoder.Encode(Metrics{ID: "NextGC", MType: "gauge", Value: &rm.NextGC}))
			check(jsonEncoder.Encode(Metrics{ID: "NumForcedGC", MType: "gauge", Value: &rm.NumForcedGC}))
			check(jsonEncoder.Encode(Metrics{ID: "NumGC", MType: "gauge", Value: &rm.NumGC}))
			check(jsonEncoder.Encode(Metrics{ID: "OtherSys", MType: "gauge", Value: &rm.OtherSys}))
			check(jsonEncoder.Encode(Metrics{ID: "PauseTotalNs", MType: "gauge", Value: &rm.PauseTotalNs}))
			check(jsonEncoder.Encode(Metrics{ID: "StackInuse", MType: "gauge", Value: &rm.StackInuse}))
			check(jsonEncoder.Encode(Metrics{ID: "StackSys", MType: "gauge", Value: &rm.StackSys}))
			check(jsonEncoder.Encode(Metrics{ID: "Sys", MType: "gauge", Value: &rm.Sys}))
			check(jsonEncoder.Encode(Metrics{ID: "RandomValue", MType: "gauge", Value: &rm.RandomValue}))
			check(jsonEncoder.Encode(Metrics{ID: "PollCount", MType: "counter", Delta: &rm.PollCount}))

			//fmt.Println(buf.Len(), buf.String())
			//hash(fmt.Sprintf("%s:counter:%d", id, delta), key)
			//hash(fmt.Sprintf("%s:gauge:%f", id, value), key)
			fmt.Printf("%s:hash:%s\n", fmt.Sprintf("%s:counter:%d", "PollCount", rm.PollCount),
				hash(fmt.Sprintf("%s:counter:%d", "PollCount", rm.PollCount), key))
			fmt.Printf("%s:hash:%s\n", fmt.Sprintf("%s:gauge:%f", "Sys", rm.Sys),
				hash(fmt.Sprintf("%s:gauge:%f", "Sys", rm.Sys), key))
			fmt.Printf("%s:hash:%s\n", fmt.Sprintf("%s:gauge:%v", "RandomValue", &rm.RandomValue),
				hash(fmt.Sprintf("%s:gauge:%f", "RandomValue", rm.RandomValue), key))
		*/
	}
}

func hash(m, k string) string {
	h := hmac.New(sha256.New, []byte(k))
	h.Write([]byte(m))
	dst := h.Sum(nil)

	//fmt.Printf("%x", dst)
	return hex.EncodeToString(dst)
}
