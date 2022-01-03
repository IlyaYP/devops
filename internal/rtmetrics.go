package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type gauge float64
type counter int64
type runTimeMetrics struct {
	Alloc,
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
type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewMonitor(buf io.Writer) func() {
	var rtm runtime.MemStats
	var rm runTimeMetrics
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
		rm.Alloc = float64(rtm.Alloc)
		rm.BuckHashSys = float64(rtm.BuckHashSys)
		rm.Frees = float64(rtm.Frees)
		rm.GCCPUFraction = rtm.GCCPUFraction
		rm.GCSys = float64(rtm.GCSys)
		rm.HeapAlloc = float64(rtm.HeapAlloc)
		rm.HeapIdle = float64(rtm.HeapIdle)
		rm.HeapInuse = float64(rtm.HeapInuse)
		rm.HeapObjects = float64(rtm.HeapObjects)
		rm.HeapReleased = float64(rtm.HeapReleased)
		rm.HeapSys = float64(rtm.HeapSys)
		rm.LastGC = float64(rtm.LastGC)
		rm.Lookups = float64(rtm.Lookups)
		rm.MCacheInuse = float64(rtm.MCacheInuse)
		rm.MCacheSys = float64(rtm.MCacheSys)
		rm.MSpanInuse = float64(rtm.MSpanInuse)
		rm.MSpanSys = float64(rtm.MSpanSys)
		rm.Mallocs = float64(rtm.Mallocs)
		rm.NextGC = float64(rtm.NextGC)
		rm.NumForcedGC = float64(rtm.NumForcedGC)
		rm.NumGC = float64(rtm.NumGC)
		rm.OtherSys = float64(rtm.OtherSys)
		rm.PauseTotalNs = float64(rtm.PauseTotalNs)
		rm.StackInuse = float64(rtm.StackInuse)
		rm.StackSys = float64(rtm.StackSys)
		rm.Sys = float64(rtm.Sys)
		rm.RandomValue = r1.Float64()
		rm.PollCount = PollCount

		//var buf bytes.Buffer
		jsonEncoder := json.NewEncoder(buf)
		check(jsonEncoder.Encode(Metric{ID: "Alloc", MType: "gauge", Value: &rm.Alloc}))
		check(jsonEncoder.Encode(Metric{ID: "BuckHashSys", MType: "gauge", Value: &rm.BuckHashSys}))
		check(jsonEncoder.Encode(Metric{ID: "Frees", MType: "gauge", Value: &rm.Frees}))
		check(jsonEncoder.Encode(Metric{ID: "GCCPUFraction", MType: "gauge", Value: &rm.GCCPUFraction}))
		check(jsonEncoder.Encode(Metric{ID: "GCSys", MType: "gauge", Value: &rm.GCSys}))
		check(jsonEncoder.Encode(Metric{ID: "HeapAlloc", MType: "gauge", Value: &rm.HeapAlloc}))
		check(jsonEncoder.Encode(Metric{ID: "HeapIdle", MType: "gauge", Value: &rm.HeapIdle}))
		check(jsonEncoder.Encode(Metric{ID: "HeapInuse", MType: "gauge", Value: &rm.HeapInuse}))
		check(jsonEncoder.Encode(Metric{ID: "HeapObjects", MType: "gauge", Value: &rm.HeapObjects}))
		check(jsonEncoder.Encode(Metric{ID: "HeapReleased", MType: "gauge", Value: &rm.HeapReleased}))
		check(jsonEncoder.Encode(Metric{ID: "HeapSys", MType: "gauge", Value: &rm.HeapSys}))
		check(jsonEncoder.Encode(Metric{ID: "LastGC", MType: "gauge", Value: &rm.LastGC}))
		check(jsonEncoder.Encode(Metric{ID: "Lookups", MType: "gauge", Value: &rm.Lookups}))
		check(jsonEncoder.Encode(Metric{ID: "MCacheInuse", MType: "gauge", Value: &rm.MCacheInuse}))
		check(jsonEncoder.Encode(Metric{ID: "MCacheSys", MType: "gauge", Value: &rm.MCacheSys}))
		check(jsonEncoder.Encode(Metric{ID: "MSpanInuse", MType: "gauge", Value: &rm.MSpanInuse}))
		check(jsonEncoder.Encode(Metric{ID: "MSpanSys", MType: "gauge", Value: &rm.MSpanSys}))
		check(jsonEncoder.Encode(Metric{ID: "Mallocs", MType: "gauge", Value: &rm.Mallocs}))
		check(jsonEncoder.Encode(Metric{ID: "NextGC", MType: "gauge", Value: &rm.NextGC}))
		check(jsonEncoder.Encode(Metric{ID: "NumForcedGC", MType: "gauge", Value: &rm.NumForcedGC}))
		check(jsonEncoder.Encode(Metric{ID: "NumGC", MType: "gauge", Value: &rm.NumGC}))
		check(jsonEncoder.Encode(Metric{ID: "OtherSys", MType: "gauge", Value: &rm.OtherSys}))
		check(jsonEncoder.Encode(Metric{ID: "PauseTotalNs", MType: "gauge", Value: &rm.PauseTotalNs}))
		check(jsonEncoder.Encode(Metric{ID: "StackInuse", MType: "gauge", Value: &rm.StackInuse}))
		check(jsonEncoder.Encode(Metric{ID: "StackSys", MType: "gauge", Value: &rm.StackSys}))
		check(jsonEncoder.Encode(Metric{ID: "Sys", MType: "gauge", Value: &rm.Sys}))
		check(jsonEncoder.Encode(Metric{ID: "RandomValue", MType: "gauge", Value: &rm.RandomValue}))
		check(jsonEncoder.Encode(Metric{ID: "PollCount", MType: "counter", Delta: &rm.PollCount}))

		//fmt.Println(buf.Len(), buf.String())
	}
}

func WriteMetric(m string) error {
	mrtm := map[string]string{
		"Alloc":         "123456789",
		"BuckHashSys":   "",
		"Frees":         "",
		"GCCPUFraction": "",
		"GCSys":         "",
		"HeapAlloc":     "",
		"HeapIdle":      "",
		"HeapInuse":     "",
		"HeapObjects":   "",
		"HeapReleased":  "",
		"HeapSys":       "",
		"LastGC":        "",
		"Lookups":       "",
		"MCacheInuse":   "",
		"MCacheSys":     "",
		"MSpanInuse":    "",
		"MSpanSys":      "",
		"Mallocs":       "",
		"NextGC":        "",
		"NumForcedGC":   "",
		"NumGC":         "",
		"OtherSys":      "",
		"PauseTotalNs":  "",
		"StackInuse":    "",
		"StackSys":      "",
		"Sys":           "",
		"RandomValue":   "",
		"PollCount":     "",
		"testCounter":   "",
	}
	k := strings.Split(m, "/")

	_, ok := mrtm[k[3]]
	if !ok {
		return fmt.Errorf("no such metric")
	}

	v, err := strconv.ParseFloat(k[4], 64)
	if err != nil {
		return err
	}
	fmt.Println("request URL:", k[3], v)

	//return nil, fmt.Errorf("orderProcessorSvc init: %w", err)
	return nil
}
