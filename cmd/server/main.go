package main

import (
	"github.com/IlyaYP/devops/cmd/server/handlers"
	"github.com/IlyaYP/devops/storage/inmemory"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	st := inmemory.NewStorage()
	r := chi.NewRouter()
	//r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
	//	rw.Write([]byte("chi"))
	//})
	r.Post("/update/{MType}/{MName}/{MVal}", handlers.UpdateHandler(st))
	http.ListenAndServe(":8080", r)

	//func(w http.ResponseWriter, r *http.Request) {
	//	handlers.UpdateHandler(w, r, st)

	/*
		// маршрутизация запросов обработчику
		//http.HandleFunc("/", handlers.HelloWorld)
		http.HandleFunc("/update/", handlers.UpdateHandler)
		// запуск сервера с адресом localhost, порт 8080
		log.Fatal(http.ListenAndServe(":8080", nil))
	*/
}
