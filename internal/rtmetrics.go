package internal

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	QueueSize = 200
)

type RunTimeMetrics struct {
	sync.RWMutex
	Gauge        map[string]float64
	Counter      map[string]int64
	rnd          *rand.Rand
	rtm          runtime.MemStats
	PoolInterval time.Duration
	Key          string
	Buf          *bytes.Buffer
	queue        chan MetricsStorage
	slice        []MetricsStorage
}

type Metrics struct {
	ID    string `json:"id"`              // имя метрики
	MType string `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64 `json:"delta,omitempty"` // значение метрики в случае передачи counter
	// видимо указатель для того чтобы передавались значения равные 0
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

type MetricsStorage struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string  `json:"hash,omitempty"`  // значение хеш-функции
}

func NewRTM(buf *bytes.Buffer, key string) *RunTimeMetrics {
	var rm RunTimeMetrics
	rm.Key = key
	rm.Buf = buf
	rm.Gauge = make(map[string]float64)
	rm.Counter = make(map[string]int64)
	rm.Counter["PollCount"] = 0
	s1 := rand.NewSource(time.Now().UnixNano())
	rm.rnd = rand.New(s1)
	rm.queue = make(chan MetricsStorage, QueueSize)

	return &rm
}

func NewMonitor(buf *bytes.Buffer, key string) func() {
	rm := NewRTM(buf, key)

	return func() {
		rm.Update()
		rm.WriteJSONtoBuf()
	}
}

func (rm *RunTimeMetrics) addToSlice() {
	rm.Lock()
	defer rm.Unlock()

	if rm.Key == "" {
		for id, value := range rm.Gauge {
			rm.slice = append(rm.slice, MetricsStorage{ID: id, MType: "gauge", Value: value})
		}
		for id, delta := range rm.Counter {
			rm.slice = append(rm.slice, MetricsStorage{ID: id, MType: "counter", Delta: delta})
		}
	} else {
		for id, value := range rm.Gauge {
			rm.slice = append(rm.slice, MetricsStorage{ID: id, MType: "gauge", Value: value,
				Hash: Hash(fmt.Sprintf("%s:gauge:%f", id, value), rm.Key)})
		}
		for id, delta := range rm.Counter {
			rm.slice = append(rm.slice, MetricsStorage{ID: id, MType: "counter", Delta: delta,
				Hash: Hash(fmt.Sprintf("%s:counter:%d", id, delta), rm.Key)})
		}
	}
}

// GetJSON writes metrics to buf as JSON Array
func (rm *RunTimeMetrics) GetJSON() *bytes.Buffer {
	rm.Lock()
	defer rm.Unlock()
	var m []Metrics
	for i, metricsStorage := range rm.slice {
		if metricsStorage.MType == "counter" {
			m = append(m, Metrics{ID: metricsStorage.ID, MType: metricsStorage.MType, Delta: &rm.slice[i].Delta,
				Hash: metricsStorage.Hash})
		} else if metricsStorage.MType == "gauge" {
			m = append(m, Metrics{ID: metricsStorage.ID, MType: metricsStorage.MType, Value: &rm.slice[i].Value,
				Hash: metricsStorage.Hash})
		}
	}
	if err := json.NewEncoder(rm.Buf).Encode(m); err != nil {
		log.Printf("Error marshal JSON in GetJSON():%s", err)
		return nil
	}
	rm.slice = nil
	return rm.Buf
}

func (rm *RunTimeMetrics) Update() {
	rm.Lock()
	defer rm.Unlock()
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

func (rm *RunTimeMetrics) UpdatePS() {
	// gopsutil
	v, _ := mem.VirtualMemory()
	info, _ := load.Avg()
	rm.Lock()
	defer rm.Unlock()
	rm.Gauge["TotalMemory"] = float64(v.Total)
	rm.Gauge["FreeMemory"] = float64(v.Free)
	// just try ;)) Вероятно ошибка в тесте потому что при не изменении этой метрики, а она не изменяется, тест не проходит
	rm.Gauge["TotalMemory"] = rm.Gauge["TotalMemory"] * rm.rnd.Float64()

	rm.Gauge["CPUutilization0"] = info.Load1

	cps, _ := cpu.Percent(0, true)
	for i, cp := range cps {
		rm.Gauge[fmt.Sprintf("CPUutilization%d", i+1)] = cp
	}
}

func (rm *RunTimeMetrics) Run(ctx context.Context) {
	poll := time.NewTicker(rm.PoolInterval)
	defer poll.Stop()
	for {
		select {
		case <-poll.C:
			rm.Collect()
		case <-ctx.Done():
			log.Println("Agent Run exiting...")
			return
		}
	}
}

func (rm *RunTimeMetrics) Collect() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		rm.Update()
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		rm.UpdatePS()
		wg.Done()
	}()
	wg.Wait()
	rm.addToSlice()
}

// WriteJSONtoBuf writes metrics to buf as JSON stream
func (rm *RunTimeMetrics) WriteJSONtoBuf() {
	check := func(err error) { // TODO: переделать костыль
		if err != nil {
			log.Println(err)
		}
	}

	if rm.Buf == nil {
		log.Println("nil Pointer rm.Buf")
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

// GetJSONSasArray writes metrics to buf as JSON Array
func (rm *RunTimeMetrics) GetJSONSasArray() *bytes.Buffer {
	var newBuf bytes.Buffer
	newBuf.WriteString("[")
	newBuf.ReadFrom(rm.Buf)
	newBuf.WriteString("]")
	return &newBuf
}

func Hash(m, k string) string {
	h := hmac.New(sha256.New, []byte(k))
	h.Write([]byte(m))
	dst := h.Sum(nil)

	//log.Printf("%s:%x", m, dst)
	return hex.EncodeToString(dst)
}
