package internal

import (
	"fmt"
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
	RandomValue gauge
	PollCount counter
}

func NewMonitor(duration int, messages chan string) {
	var rtm runtime.MemStats
	//var rm runTimeMetrics
	PollCount := 0
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	var interval = time.Duration(duration) * time.Second
	for {
		<-time.After(interval)
		runtime.ReadMemStats(&rtm)

		messages <- fmt.Sprintf("/gauge/Alloc/%v", rtm.Alloc)
		messages <- fmt.Sprintf("/gauge/BuckHashSys/%v", rtm.BuckHashSys)
		messages <- fmt.Sprintf("/gauge/Frees/%v", rtm.Frees)
		messages <- fmt.Sprintf("/gauge/GCCPUFraction/%v", rtm.GCCPUFraction)
		messages <- fmt.Sprintf("/gauge/GCSys/%v", rtm.GCSys)
		messages <- fmt.Sprintf("/gauge/HeapAlloc/%v", rtm.HeapAlloc)
		messages <- fmt.Sprintf("/gauge/HeapIdle/%v", rtm.HeapIdle)
		messages <- fmt.Sprintf("/gauge/HeapInuse/%v", rtm.HeapInuse)
		messages <- fmt.Sprintf("/gauge/HeapObjects/%v", rtm.HeapObjects)
		messages <- fmt.Sprintf("/gauge/HeapReleased/%v", rtm.HeapReleased)
		messages <- fmt.Sprintf("/gauge/HeapSys/%v", rtm.HeapSys)
		messages <- fmt.Sprintf("/gauge/LastGC/%v", rtm.LastGC)
		messages <- fmt.Sprintf("/gauge/Lookups/%v", rtm.Lookups)
		messages <- fmt.Sprintf("/gauge/MCacheInuse/%v", rtm.MCacheInuse)
		messages <- fmt.Sprintf("/gauge/MCacheSys/%v", rtm.MCacheSys)
		messages <- fmt.Sprintf("/gauge/MSpanInuse/%v", rtm.MSpanInuse)
		messages <- fmt.Sprintf("/gauge/MSpanSys/%v", rtm.MSpanSys)
		messages <- fmt.Sprintf("/gauge/Mallocs/%v", rtm.Mallocs)
		messages <- fmt.Sprintf("/gauge/NextGC/%v", rtm.NextGC)
		messages <- fmt.Sprintf("/gauge/NumForcedGC/%v", rtm.NumForcedGC)
		messages <- fmt.Sprintf("/gauge/NumGC/%v", rtm.NumGC)
		messages <- fmt.Sprintf("/gauge/OtherSys/%v", rtm.OtherSys)
		messages <- fmt.Sprintf("/gauge/PauseTotalNs/%v", rtm.PauseTotalNs)
		messages <- fmt.Sprintf("/gauge/StackInuse/%v", rtm.StackInuse)
		messages <- fmt.Sprintf("/gauge/StackSys/%v", rtm.StackSys)
		messages <- fmt.Sprintf("/gauge/Sys/%v", rtm.Sys)
		messages <- fmt.Sprintf("/gauge/RandomValue/%v", r1.Float64())
		messages <- fmt.Sprintf("/counter/PollCount/%v", PollCount)

		//fields := reflect.TypeOf(rm)
		//values := reflect.ValueOf(rm)
		//num := fields.NumField()
		//for i := 0; i < num; i++ {
		//	field := fields.Field(i)
		//	value := values.Field(i)
		//	fmt.Print("Type:", field.Type, ",", field.Name, "=", value, "\n")
		//}
		//b, _ := json.Marshal(rm)
		//fmt.Println(string(b))
		//fmt.Printf("...%v", PollCount)
		PollCount++
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
