package handlers

import (
	"fmt"
	"github.com/IlyaYP/devops/storage/inmemory"
	"net/http"
	"strconv"
	"strings"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("<h1>Hello, World</h1>"))

}

//http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
//request URL: /update/counter/PollCount/2
//request URL: /update/gauage/Alloc/201456
func UpdateHandler(st *inmemory.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println(r.URL.String())
		//fmt.Println(chi.URLParam(r, "MType"), chi.URLParam(r, "MName"),
		//	chi.URLParam(r, "MVal")) //{MType}/{MName}/{MVal}
		////update/unknown/testCounter/100
		mtypes := map[string]string{
			"gauge":   "",
			"counter": "",
		}

		//mrtm := map[string]string{
		//	"Alloc":         "123456789",
		//	"BuckHashSys":   "",
		//	"Frees":         "",
		//	"GCCPUFraction": "",
		//	"GCSys":         "",
		//	"HeapAlloc":     "",
		//	"HeapIdle":      "",
		//	"HeapInuse":     "",
		//	"HeapObjects":   "",
		//	"HeapReleased":  "",
		//	"HeapSys":       "",
		//	"LastGC":        "",
		//	"Lookups":       "",
		//	"MCacheInuse":   "",
		//	"MCacheSys":     "",
		//	"MSpanInuse":    "",
		//	"MSpanSys":      "",
		//	"Mallocs":       "",
		//	"NextGC":        "",
		//	"NumForcedGC":   "",
		//	"NumGC":         "",
		//	"OtherSys":      "",
		//	"PauseTotalNs":  "",
		//	"StackInuse":    "",
		//	"StackSys":      "",
		//	"Sys":           "",
		//	"RandomValue":   "",
		//	"PollCount":     "",
		//	"testCounter":   "",
		//	"testGauge":     "",
		//}
		k := strings.Split(r.URL.String(), "/")

		// Check Metric Type
		_, ok := mtypes[k[2]]
		if !ok {
			//http.Error(w, "no such type", http.StatusInternalServerError)
			http.Error(w, "no such type", http.StatusNotImplemented)
			return
		}

		//Check Metric Name
		//_, ok = mrtm[k[3]]
		//if !ok {
		//	http.Error(w, "no such metric", http.StatusNotFound)
		//	return
		//}

		// Check Metric Value
		v, err := strconv.ParseFloat(k[4], 64)
		if err != nil {
			http.Error(w, "not found", http.StatusBadRequest)
			return
		}
		fmt.Println("request URL:", k[3], v)

		//err := internal.WriteMetric(r.URL.String())
		//if err != nil {
		//	http.Error(w, "not found", http.StatusNotFound)
		//	return
		//}
		////fmt.Println("request Headers:", r.Header)
		//body, _ := io.ReadAll(r.Body)
		//fmt.Println("request Body:", string(body))

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		//_, err = w.Write([]byte("OK"))
		//if err != nil {
		//	return
		//}
	}
}

//var mux map[string]int
//m := make(map[string]int)
