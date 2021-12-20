package main

import (
	"fmt"
	"net/http"
)

func update(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request URL:", r.URL)
	////fmt.Println("request Headers:", r.Header)
	//body, _ := io.ReadAll(r.Body)
	//fmt.Println("request Body:", string(body))

	_, err := w.Write([]byte("OK"))
	if err != nil {
		return 
	}
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello, World</h1>"))
}

func main() {
    // маршрутизация запросов обработчику
    http.HandleFunc("/", HelloWorld)
	http.HandleFunc("/update/", update)
    // запуск сервера с адресом localhost, порт 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return 
	}
} 




/*
http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
*/