package internal

import (
	"bytes"
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

func NewRTM() *RunTimeMetrics {
	var rm RunTimeMetrics
	rm.Gauge = make(map[string]float64)
	rm.Counter = make(map[string]int64)
	rm.Counter["PollCount"] = 0
	s1 := rand.NewSource(time.Now().UnixNano())
	rm.rnd = rand.New(s1)
	return &rm
}

type RunTimeMetrics struct {
	Gauge        map[string]float64
	Counter      map[string]int64
	rnd          *rand.Rand
	rtm          runtime.MemStats
	PoolInterval time.Duration
	Key          string
	Buf          *bytes.Buffer
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

type MetricStorage struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string  `json:"hash,omitempty"`  // значение хеш-функции
}

func (rm *RunTimeMetrics) Update() {
	runtime.ReadMemStats(&rm.rtm)

	rm.Counter["PollCount"]++
	rm.Gauge["RandomValue"] = rm.rnd.Float64()
	rm.Gauge["Alloc"] = float64(rm.rtm.Alloc)
	rm.Gauge["TotalAlloc"] = float64(rm.rtm.TotalAlloc)
	rm.Gauge["BuckHashSys"] = float64(rm.rtm.BuckHashSys)
	rm.Gauge["Frees"] = float64(rm.rtm.Frees)
	rm.Gauge["GCCPUFraction"] = rm.rtm.GCCPUFraction
	rm.Gauge["GCSys"] = float64(rm.rtm.GCSys)
	rm.Gauge["HeapAlloc"] = float64(rm.rtm.HeapAlloc)
	rm.Gauge["HeapIdle"] = float64(rm.rtm.HeapIdle)
	rm.Gauge["HeapInuse"] = float64(rm.rtm.HeapInuse)
	rm.Gauge["HeapObjects"] = float64(rm.rtm.HeapObjects)
	rm.Gauge["HeapReleased"] = float64(rm.rtm.HeapReleased)
	rm.Gauge["HeapSys"] = float64(rm.rtm.HeapSys)
	rm.Gauge["LastGC"] = float64(rm.rtm.LastGC)
	rm.Gauge["Lookups"] = float64(rm.rtm.Lookups)
	rm.Gauge["MCacheInuse"] = float64(rm.rtm.MCacheInuse)
	rm.Gauge["MCacheSys"] = float64(rm.rtm.MCacheSys)
	rm.Gauge["MSpanInuse"] = float64(rm.rtm.MSpanInuse)
	rm.Gauge["MSpanSys"] = float64(rm.rtm.MSpanSys)
	rm.Gauge["Mallocs"] = float64(rm.rtm.Mallocs)
	rm.Gauge["NextGC"] = float64(rm.rtm.NextGC)
	rm.Gauge["NumForcedGC"] = float64(rm.rtm.NumForcedGC)
	rm.Gauge["NumGC"] = float64(rm.rtm.NumGC)
	rm.Gauge["OtherSys"] = float64(rm.rtm.OtherSys)
	rm.Gauge["PauseTotalNs"] = float64(rm.rtm.PauseTotalNs)
	rm.Gauge["StackInuse"] = float64(rm.rtm.StackInuse)
	rm.Gauge["StackSys"] = float64(rm.rtm.StackSys)
	rm.Gauge["Sys"] = float64(rm.rtm.Sys)
}

func (rm *RunTimeMetrics) Run() {
	poll := time.Tick(rm.PoolInterval)
	for {
		<-poll
		rm.Update()
	}
}

// GetJSON writes metrics to buf as JSON stream
func (rm *RunTimeMetrics) GetJSON() {
	check := func(err error) { // TODO: переделать костыль
		if err != nil {
			log.Println(err)
		}
	}

	if rm.Buf == nil {
		log.Println("nil Pointer rm.Buff")
		return
	}
	jsonEncoder := json.NewEncoder(rm.Buf)
	if rm.Key == "" {
		for id, value := range rm.Gauge {
			check(jsonEncoder.Encode(Metrics{ID: id, MType: "gauge", Value: &value}))
		}
		for id, delta := range rm.Counter {
			check(jsonEncoder.Encode(Metrics{ID: id, MType: "counter", Delta: &delta}))
		}
	} else {
		for id, value := range rm.Gauge {
			check(jsonEncoder.Encode(Metrics{ID: id, MType: "gauge", Value: &value,
				Hash: Hash(fmt.Sprintf("%s:gauge:%f", id, value), rm.Key)}))
		}
		for id, delta := range rm.Counter {
			check(jsonEncoder.Encode(Metrics{ID: id, MType: "counter", Delta: &delta,
				Hash: Hash(fmt.Sprintf("%s:counter:%d", id, delta), rm.Key)}))
		}
	}
}

func getMetrics() []MetricStorage {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	var m []MetricStorage

	m = append(m, MetricStorage{ID: "Alloc", MType: "gauge", Value: float64(rtm.Alloc)})
	m = append(m, MetricStorage{ID: "TotalAlloc", MType: "gauge", Value: float64(rtm.TotalAlloc)})
	m = append(m, MetricStorage{ID: "BuckHashSys", MType: "gauge", Value: float64(rtm.BuckHashSys)})
	m = append(m, MetricStorage{ID: "Frees", MType: "gauge", Value: float64(rtm.Frees)})
	m = append(m, MetricStorage{ID: "GCCPUFraction", MType: "gauge", Value: rtm.GCCPUFraction})
	m = append(m, MetricStorage{ID: "GCSys", MType: "gauge", Value: float64(rtm.GCSys)})
	m = append(m, MetricStorage{ID: "HeapAlloc", MType: "gauge", Value: float64(rtm.HeapAlloc)})
	m = append(m, MetricStorage{ID: "HeapIdle", MType: "gauge", Value: float64(rtm.HeapIdle)})
	m = append(m, MetricStorage{ID: "HeapInuse", MType: "gauge", Value: float64(rtm.HeapInuse)})
	m = append(m, MetricStorage{ID: "HeapObjects", MType: "gauge", Value: float64(rtm.HeapObjects)})
	m = append(m, MetricStorage{ID: "HeapReleased", MType: "gauge", Value: float64(rtm.HeapReleased)})
	m = append(m, MetricStorage{ID: "HeapSys", MType: "gauge", Value: float64(rtm.HeapSys)})
	m = append(m, MetricStorage{ID: "LastGC", MType: "gauge", Value: float64(rtm.LastGC)})
	m = append(m, MetricStorage{ID: "Lookups", MType: "gauge", Value: float64(rtm.Lookups)})
	m = append(m, MetricStorage{ID: "MCacheInuse", MType: "gauge", Value: float64(rtm.MCacheInuse)})
	m = append(m, MetricStorage{ID: "MCacheSys", MType: "gauge", Value: float64(rtm.MCacheSys)})
	m = append(m, MetricStorage{ID: "MSpanInuse", MType: "gauge", Value: float64(rtm.MSpanInuse)})
	m = append(m, MetricStorage{ID: "MSpanSys", MType: "gauge", Value: float64(rtm.MSpanSys)})
	m = append(m, MetricStorage{ID: "Mallocs", MType: "gauge", Value: float64(rtm.Mallocs)})
	m = append(m, MetricStorage{ID: "NextGC", MType: "gauge", Value: float64(rtm.NextGC)})
	m = append(m, MetricStorage{ID: "NumForcedGC", MType: "gauge", Value: float64(rtm.NumForcedGC)})
	m = append(m, MetricStorage{ID: "NumGC", MType: "gauge", Value: float64(rtm.NumGC)})
	m = append(m, MetricStorage{ID: "OtherSys", MType: "gauge", Value: float64(rtm.OtherSys)})
	m = append(m, MetricStorage{ID: "PauseTotalNs", MType: "gauge", Value: float64(rtm.PauseTotalNs)})
	m = append(m, MetricStorage{ID: "StackInuse", MType: "gauge", Value: float64(rtm.StackInuse)})
	m = append(m, MetricStorage{ID: "StackSys", MType: "gauge", Value: float64(rtm.StackSys)})
	m = append(m, MetricStorage{ID: "Sys", MType: "gauge", Value: float64(rtm.Sys)})
	return m
}

func NewMonitor(buf io.Writer, key string) func() {
	//	rm :=NewRTM()

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
					Hash: Hash(fmt.Sprintf("%s:gauge:%f", id, value), key)}))
			}
			for id, delta := range rm.Counter {
				check(jsonEncoder.Encode(Metrics{ID: id, MType: "counter", Delta: &delta,
					Hash: Hash(fmt.Sprintf("%s:counter:%d", id, delta), key)}))
			}
		}
	}
}

func Hash(m, k string) string {
	h := hmac.New(sha256.New, []byte(k))
	h.Write([]byte(m))
	dst := h.Sum(nil)

	//log.Printf("%s:%x", m, dst)
	return hex.EncodeToString(dst)
}
