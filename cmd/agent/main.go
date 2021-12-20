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
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

func send(endpoint string) {
	// Имя метрики: "Alloc", тип: gauge
	// http://localhost:8080/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	//endpoint := "http://localhost:8080/update/gauge/Alloc/4000000000.0001"
	//endpoint := "http://localhost:8080/asedffggf"
	data := `Hi, how are you? русский текст`
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(data))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	request.Header.Set("Content-Type", "text/plain")
	request.Header.Set("Content-Length", strconv.Itoa(len(data)))
	request.Header.Set("application-type", "text/plain")
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(body))
}

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

		messages<- fmt.Sprintf("/gauage/Alloc/%v", rtm.Alloc)
		messages<- fmt.Sprintf("/gauage/BuckHashSys/%v", rtm.BuckHashSys)
		messages<- fmt.Sprintf("/gauage/Frees/%v", rtm.Frees)
		messages<- fmt.Sprintf("/gauage/GCCPUFraction/%v", rtm.GCCPUFraction)
		messages<- fmt.Sprintf("/gauage/GCSys/%v", rtm.GCSys)
		messages<- fmt.Sprintf("/gauage/HeapAlloc/%v", rtm.HeapAlloc)
		messages<- fmt.Sprintf("/gauage/HeapIdle/%v", rtm.HeapIdle)
		messages<- fmt.Sprintf("/gauage/HeapInuse/%v", rtm.HeapInuse)
		messages<- fmt.Sprintf("/gauage/HeapObjects/%v", rtm.HeapObjects)
		messages<- fmt.Sprintf("/gauage/HeapReleased/%v", rtm.HeapReleased)
		messages<- fmt.Sprintf("/gauage/HeapSys/%v", rtm.HeapSys)
		messages<- fmt.Sprintf("/gauage/LastGC/%v", rtm.LastGC)
		messages<- fmt.Sprintf("/gauage/Lookups/%v", rtm.Lookups)
		messages<- fmt.Sprintf("/gauage/MCacheInuse/%v", rtm.MCacheInuse)
		messages<- fmt.Sprintf("/gauage/MCacheSys/%v", rtm.MCacheSys)
		messages<- fmt.Sprintf("/gauage/MSpanInuse/%v", rtm.MSpanInuse)
		messages<- fmt.Sprintf("/gauage/MSpanSys/%v", rtm.MSpanSys)
		messages<- fmt.Sprintf("/gauage/Mallocs/%v", rtm.Mallocs)
		messages<- fmt.Sprintf("/gauage/NextGC/%v", rtm.NextGC)
		messages<- fmt.Sprintf("/gauage/NumForcedGC/%v", rtm.NumForcedGC)
		messages<- fmt.Sprintf("/gauage/NumGC/%v", rtm.NumGC)
		messages<- fmt.Sprintf("/gauage/OtherSys/%v", rtm.OtherSys)
		messages<- fmt.Sprintf("/gauage/PauseTotalNs/%v", rtm.PauseTotalNs)
		messages<- fmt.Sprintf("/gauage/StackInuse/%v", rtm.StackInuse)
		messages<- fmt.Sprintf("/gauage/StackSys/%v", rtm.StackSys)
		messages<- fmt.Sprintf("/gauage/Sys/%v", rtm.Sys)
		messages<- fmt.Sprintf("/gauage/RandomValue/%v", r1.Float64())
		messages<- fmt.Sprintf("/counter/PollCount/%v", PollCount)


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
func main() {
	//send()
	messages := make(chan string, 200)

	go NewMonitor(2, messages)

	for{
		select {
		case str := <-messages:
			send("http://localhost:8080/update" + str)
		default:
			time.Sleep(10 * time.Second)
		}
	}

	//var input string
	//fmt.Scanln(&input)
}
