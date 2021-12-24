package main

import (
	"github.com/IlyaYP/devops/cmd/server/handlers"
	"net/http"
)



func main() {
    // маршрутизация запросов обработчику
    http.HandleFunc("/", handlers.HelloWorld)
	http.HandleFunc("/update/", handlers.UpdateHandler)
    // запуск сервера с адресом localhost, порт 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return 
	}
}

