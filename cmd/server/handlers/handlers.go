package handlers

import (
	"context"
	"github.com/IlyaYP/devops/storage/inmemory"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
)

func ReadHandler(st *inmemory.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := st.ReadMetrics(context.Background())
		w.WriteHeader(http.StatusOK)
		var metrics []string
		for k, v := range m {
			for kk, vv := range v {
				metrics = append(metrics, k+" "+kk+": "+vv)
			}
		}
		sort.Strings(metrics)

		const tpl = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
  </head>
  <body>
    {{range .Items}}<div>{{.}}</div>{{else}}<div><strong>no data</strong></div>{{end}}
  </body>
</html>`
		check := func(err error) {
			if err != nil {
				log.Fatal(err)
			}
		}
		t, err := template.New("webpage").Parse(tpl)
		check(err)
		data := struct {
			Title string
			Items []string
		}{
			Title: "Metrics",
			Items: metrics,
		}

		err = t.Execute(w, data)
		check(err)
	}
}

//GET http://localhost:8080/value/counter/testSetGet33
//GET http://localhost:8080/value/counter/PollCount
func GetHandler(st *inmemory.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		k := strings.Split(r.URL.String(), "/") // TODO: Chi not work in tests, so using old method
		//if err := st.PutMetric(context.Background(), chi.URLParam(r, "MType"),
		//	chi.URLParam(r, "MName"), chi.URLParam(r, "MVal")); err != nil {
		v, err := st.GetMetric(context.Background(), k[2], k[3])
		if err != nil {
			if err.Error() == "wrong type" {
				http.Error(w, err.Error(), http.StatusNotImplemented)
			} else if err.Error() == "no such metric" {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte(v)); err != nil {
			return
		}

	}
}

//http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
//request URL: /update/counter/PollCount/2
//request URL: /update/gauage/Alloc/201456
func UpdateHandler(st *inmemory.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		k := strings.Split(r.URL.String(), "/") // TODO: Chi not work in tests, so using old method
		//if err := st.PutMetric(context.Background(), chi.URLParam(r, "MType"),
		//	chi.URLParam(r, "MName"), chi.URLParam(r, "MVal")); err != nil {
		if err := st.PutMetric(context.Background(), k[2], k[3], k[4]); err != nil {
			if err.Error() == "wrong type" {
				http.Error(w, err.Error(), http.StatusNotImplemented)
			} else if err.Error() == "wrong value" {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(w, "unknown error", http.StatusBadRequest)
			}
			return
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
