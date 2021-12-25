package handlers

import (
	"fmt"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("<h1>Hello, World</h1>"))

}

//http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
//request URL: /update/counter/PollCount/2
//request URL: /update/gauage/Alloc/201456
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request URL:", r.URL)
	////fmt.Println("request Headers:", r.Header)
	//body, _ := io.ReadAll(r.Body)
	//fmt.Println("request Body:", string(body))

	//user, ok := Metrics[Metric]
	//if !ok {
	//	http.Error(r, "user not found", http.StatusNotFound)
	//	return
	//}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}

//var mux map[string]int
//m := make(map[string]int)
