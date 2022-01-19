package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/IlyaYP/devops/internal"
	"github.com/IlyaYP/devops/storage"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Handlers struct {
	St storage.MetricStorage // Q: почему сюда не могу поставить указатель? не создаю ли я копию?
}

func (h *Handlers) ReadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := h.St.ReadMetrics()
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		//w.Header().Set("Content-Encoding", "gzip")
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

// GetHandler receiving requests like these, and responds value in body
//GET http://localhost:8080/value/counter/testSetGet33
//GET http://localhost:8080/value/counter/PollCount
func (h *Handlers) GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		k := strings.Split(r.URL.String(), "/") // TODO: Chi not work in tests, so using old method
		//if err := st.PutMetric(context.Background(), chi.URLParam(r, "MType"),
		//	chi.URLParam(r, "MName"), chi.URLParam(r, "MVal")); err != nil {
		v, err := h.St.GetMetric(k[2], k[3])
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
			log.Println("GetHandler Write body:", err)
			return
		}

	}
}

// UpdateHandler serves following requests:
//http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
//request URL: /update/counter/PollCount/2
//request URL: /update/gauage/Alloc/201456
func (h *Handlers) UpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		k := strings.Split(r.URL.String(), "/") // TODO: Chi not work in tests, so using old method
		//if err := st.PutMetric(context.Background(), chi.URLParam(r, "MType"),
		//	chi.URLParam(r, "MName"), chi.URLParam(r, "MVal")); err != nil {
		if err := h.St.PutMetric(k[2], k[3], k[4]); err != nil {
			if err.Error() == "wrong type" {
				http.Error(w, err.Error(), http.StatusNotImplemented)
			} else if err.Error() == "wrong value" {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(w, "unknown error", http.StatusBadRequest)
			}
			log.Println("UpdateHandler:", err)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

// UpdateJSONHandler receiving updates in JSON in body
//POST http://localhost:8080/update/
func (h *Handlers) UpdateJSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonDecoder := json.NewDecoder(r.Body)
		for jsonDecoder.More() {
			var m internal.Metrics
			var MetricValue string

			err := jsonDecoder.Decode(&m)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Println("UpdateJSONHandler:jsonDecoder.Decode", err)
				return
			}

			if m.MType == "gauge" && m.Value != nil {
				MetricValue = fmt.Sprintf("%v", *m.Value)
			} else if m.MType == "counter" && m.Delta != nil {
				MetricValue = fmt.Sprintf("%v", *m.Delta)
			} else {
				http.Error(w, "wrong type", http.StatusNotImplemented)
				log.Println("UpdateJSONHandler:", err)
				return
			}

			if err := h.St.PutMetric(m.MType, m.ID, MetricValue); err != nil {
				if err.Error() == "wrong type" {
					http.Error(w, err.Error(), http.StatusNotImplemented)
				} else if err.Error() == "wrong value" {
					http.Error(w, err.Error(), http.StatusBadRequest)
				} else {
					http.Error(w, "unknown error", http.StatusBadRequest)
				}
				log.Println("UpdateJSONHandler:PutMetric ", err)
				return
			}

		}

		w.Header().Set("content-type", "application/json")
		w.Header().Set("application-type", "text/plain")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil { // no need but doesn't pass test without it
			log.Println("UpdateJSONHandler Write body:", err)
			return
		}
	}
}

// GetJSONHandler receiving requests in JSON body, and responds via JSON in body
//POST http://localhost:8080/value/
func (h *Handlers) GetJSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonDecoder := json.NewDecoder(r.Body)
		jsonEncoder := json.NewEncoder(w)

		w.Header().Set("content-type", "application/json")
		//w.WriteHeader(http.StatusOK)

		// while the r.body  contains values
		for jsonDecoder.More() {
			var m internal.Metrics

			err := jsonDecoder.Decode(&m)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Println("GetJSONHandler:jsonDecoder.Decode:185", err)
				return
			}
			if m.MType != "gauge" && m.MType != "counter" {
				http.Error(w, "wrong type", http.StatusNotImplemented)
				log.Println("GetJSONHandler:190", err)
				return
			}

			v, err := h.St.GetMetric(m.MType, m.ID)
			if err != nil {
				if err.Error() == "wrong type" {
					http.Error(w, err.Error(), http.StatusNotImplemented)
				} else if err.Error() == "no such metric" {
					http.Error(w, err.Error(), http.StatusNotFound)
				} else {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				log.Println("GetJSONHandler:GetMetric:", err)
				return
			}
			if m.MType == "gauge" {
				value, err := strconv.ParseFloat(v, 64)
				if err != nil {
					log.Println(err)
				}
				m.Value = &value
			} else if m.MType == "counter" {
				delta, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					log.Println(err)
				}
				m.Delta = &delta
			}
			//w.Header().Set("content-type", "application/json")
			//w.WriteHeader(http.StatusOK)
			//log.Println("RESP:", m.MType, m.ID, v) // DEBUG:

			if err := jsonEncoder.Encode(m); err != nil {
				log.Println("GetJSONHandler write JSON body", err)
			}
		}
	}
}
