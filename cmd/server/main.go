package main

import (
	"context"
	"fmt"
	"github.com/IlyaYP/devops/cmd/server/handlers"
	"github.com/IlyaYP/devops/storage/inmemory"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func testStore(st *inmemory.Storage) {
	println(st)

	if err := st.PutMetric(context.Background(), "gauge", "aaa", "111.111"); err != nil {
		println(err.Error())
		return
	}

	if err := st.PutMetric(context.Background(), "counter", "bbb", "222"); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := st.PutMetric(context.Background(), "gauge", "ccc", "333.333"); err != nil {
		fmt.Println(err.Error())
		return
	}

	if v, err := st.GetMetric(context.Background(), "gauge", "ccc"); err != nil {
		fmt.Println(err.Error())
		return
	} else {
		fmt.Println(v)
	}
	ret := st.ReadMetrics(context.Background())
	fmt.Println(ret)

}

func main() {
	st := inmemory.NewStorage()
	//testStore(st)

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
