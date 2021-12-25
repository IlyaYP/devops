package main

import (
	"github.com/IlyaYP/devops/cmd/server/handlers"
	"log"
	"net/http"
)

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", handlers.HelloWorld)
	http.HandleFunc("/update/", handlers.UpdateHandler)
	// запуск сервера с адресом localhost, порт 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
