package internal

import (
	"fmt"
	"math/rand"
	"runtime"
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

		messages <- fmt.Sprintf("/gauage/Alloc/%v", rtm.Alloc)
		messages <- fmt.Sprintf("/gauage/BuckHashSys/%v", rtm.BuckHashSys)
		messages <- fmt.Sprintf("/gauage/Frees/%v", rtm.Frees)
		messages <- fmt.Sprintf("/gauage/GCCPUFraction/%v", rtm.GCCPUFraction)
		messages <- fmt.Sprintf("/gauage/GCSys/%v", rtm.GCSys)
		messages <- fmt.Sprintf("/gauage/HeapAlloc/%v", rtm.HeapAlloc)
		messages <- fmt.Sprintf("/gauage/HeapIdle/%v", rtm.HeapIdle)
		messages <- fmt.Sprintf("/gauage/HeapInuse/%v", rtm.HeapInuse)
		messages <- fmt.Sprintf("/gauage/HeapObjects/%v", rtm.HeapObjects)
		messages <- fmt.Sprintf("/gauage/HeapReleased/%v", rtm.HeapReleased)
		messages <- fmt.Sprintf("/gauage/HeapSys/%v", rtm.HeapSys)
		messages <- fmt.Sprintf("/gauage/LastGC/%v", rtm.LastGC)
		messages <- fmt.Sprintf("/gauage/Lookups/%v", rtm.Lookups)
		messages <- fmt.Sprintf("/gauage/MCacheInuse/%v", rtm.MCacheInuse)
		messages <- fmt.Sprintf("/gauage/MCacheSys/%v", rtm.MCacheSys)
		messages <- fmt.Sprintf("/gauage/MSpanInuse/%v", rtm.MSpanInuse)
		messages <- fmt.Sprintf("/gauage/MSpanSys/%v", rtm.MSpanSys)
		messages <- fmt.Sprintf("/gauage/Mallocs/%v", rtm.Mallocs)
		messages <- fmt.Sprintf("/gauage/NextGC/%v", rtm.NextGC)
		messages <- fmt.Sprintf("/gauage/NumForcedGC/%v", rtm.NumForcedGC)
		messages <- fmt.Sprintf("/gauage/NumGC/%v", rtm.NumGC)
		messages <- fmt.Sprintf("/gauage/OtherSys/%v", rtm.OtherSys)
		messages <- fmt.Sprintf("/gauage/PauseTotalNs/%v", rtm.PauseTotalNs)
		messages <- fmt.Sprintf("/gauage/StackInuse/%v", rtm.StackInuse)
		messages <- fmt.Sprintf("/gauage/StackSys/%v", rtm.StackSys)
		messages <- fmt.Sprintf("/gauage/Sys/%v", rtm.Sys)
		messages <- fmt.Sprintf("/gauage/RandomValue/%v", r1.Float64())
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
	fmt.Println("request URL:", k[3])

	_, ok := mrtm[k[3]]
	if !ok {
		return fmt.Errorf("no such metric")
	}

	//return nil, fmt.Errorf("orderProcessorSvc init: %w", err)
	return nil
}
