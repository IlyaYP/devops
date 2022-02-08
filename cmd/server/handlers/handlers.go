package handlers

import (
	"bytes"
	"compress/flate"
	"crypto/hmac"
	"encoding/json"
	"fmt"
	"github.com/IlyaYP/devops/internal"
	"github.com/IlyaYP/devops/storage"
	"html/template"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Handlers struct {
	St  storage.MetricStorage
	Key string
}

func (h *Handlers) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if err := h.St.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func (h *Handlers) ReadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := h.St.ReadMetrics(r.Context())
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
		v, err := h.St.GetMetric(r.Context(), k[2], k[3])
		if err != nil {
			switch err.(type) {
			case *storage.TypeError:
				http.Error(w, err.Error(), http.StatusNotImplemented)
			case *storage.MetricError:
				http.Error(w, err.Error(), http.StatusNotFound)
			default:
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			//if err.Error() == "wrong type" {
			//	http.Error(w, err.Error(), http.StatusNotImplemented)
			//} else if strings.HasPrefix(err.Error(), "no such metric") {
			//	http.Error(w, err.Error(), http.StatusNotFound)
			//} else {
			//	http.Error(w, err.Error(), http.StatusBadRequest)
			//}
			log.Println(err)
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
//request URL: /update/gauge/Alloc/201456
func (h *Handlers) UpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		k := strings.Split(r.URL.String(), "/") // TODO: Chi not work in tests, so using old method
		//if err := st.PutMetric(context.Background(), chi.URLParam(r, "MType"),
		//	chi.URLParam(r, "MName"), chi.URLParam(r, "MVal")); err != nil {
		if err := h.St.PutMetric(r.Context(), k[2], k[3], k[4]); err != nil {
			//if err.Error() == "wrong type" {
			//	http.Error(w, err.Error(), http.StatusNotImplemented)
			//} else if err.Error() == "wrong value" {
			//	http.Error(w, err.Error(), http.StatusBadRequest)
			//} else {
			//	http.Error(w, "unknown error", http.StatusBadRequest)
			//}
			switch err.(type) {
			case *storage.TypeError:
				http.Error(w, err.Error(), http.StatusNotImplemented)
			case *storage.MetricError:
				http.Error(w, err.Error(), http.StatusNotFound)
			default:
				http.Error(w, err.Error(), http.StatusBadRequest)
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
				if h.Key != "" {
					hash := internal.Hash(fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value), h.Key)
					if !hmac.Equal([]byte(hash), []byte(m.Hash)) {
						http.Error(w, "wrong hash", http.StatusBadRequest)
						log.Printf("Wrong Hash of %s:gauge:%f\n%s\n%s", m.ID, *m.Value,
							m.Hash, hash)
						return
					}
				}

			} else if m.MType == "counter" && m.Delta != nil {
				MetricValue = fmt.Sprintf("%v", *m.Delta)
				if h.Key != "" {
					hash := internal.Hash(fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta), h.Key)
					if !hmac.Equal([]byte(hash), []byte(m.Hash)) {
						http.Error(w, "wrong hash", http.StatusBadRequest)
						log.Printf("Wrong Hash of %s:counter:%d\n%s\n%s", m.ID, *m.Delta,
							m.Hash, hash)
						return
					}
				}
			} else {
				http.Error(w, "wrong type", http.StatusNotImplemented)
				log.Println("UpdateJSONHandler:", m)
				return
			}

			if err := h.St.PutMetric(r.Context(), m.MType, m.ID, MetricValue); err != nil {
				//if err.Error() == "wrong type" {
				//	http.Error(w, err.Error(), http.StatusNotImplemented)
				//} else if err.Error() == "wrong value" {
				//	http.Error(w, err.Error(), http.StatusBadRequest)
				//} else {
				//	http.Error(w, "unknown error", http.StatusBadRequest)
				//}
				switch err.(type) {
				case *storage.TypeError:
					http.Error(w, err.Error(), http.StatusNotImplemented)
				default:
					http.Error(w, err.Error(), http.StatusBadRequest)
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

// UpdatesJSONHandler receiving updates in JSON array in body
//POST http://localhost:8080/updates/
func (h *Handlers) UpdatesJSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//log.Println("UpdatesJSONHandler:", r.Header.Values("Content-Encoding"))
		//log.Println("UpdatesJSONHandler:", r.Header.Values("Content-Length"))

		enc := r.Header.Values("Content-Encoding")
		var gzip bool
		for _, v := range enc {
			if v == "gzip" {
				gzip = true
			}
		}
		var jsonDecoder *json.Decoder
		if gzip {
			body, _ := io.ReadAll(r.Body)
			defer r.Body.Close()

			rr := flate.NewReader(bytes.NewReader(body))
			defer rr.Close()

			var b bytes.Buffer
			// в переменную b записываются распакованные данные
			_, err := b.ReadFrom(rr)
			if err != nil {
				log.Printf("failed decompress data: %v", err)
				return
			}

			//log.Println("request Body:", b.String())
			jsonDecoder = json.NewDecoder(&b)
		} else {
			jsonDecoder = json.NewDecoder(r.Body)
		}
		for jsonDecoder.More() {
			var mm []internal.Metrics
			var MetricValue string

			err := jsonDecoder.Decode(&mm)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Println("UpdatesJSONHandler:jsonDecoder.Decode", err)
				return
			}

			for _, m := range mm {

				if m.MType == "gauge" && m.Value != nil {
					MetricValue = fmt.Sprintf("%v", *m.Value)
					if h.Key != "" {
						hash := internal.Hash(fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value), h.Key)
						if !hmac.Equal([]byte(hash), []byte(m.Hash)) {
							http.Error(w, "wrong hash", http.StatusBadRequest)
							log.Printf("Wrong Hash of %s:gauge:%f\n%s\n%s", m.ID, *m.Value,
								m.Hash, hash)
							return
						}
					}
				} else if m.MType == "counter" && m.Delta != nil {
					MetricValue = fmt.Sprintf("%v", *m.Delta)
					if h.Key != "" {
						hash := internal.Hash(fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta), h.Key)
						if !hmac.Equal([]byte(hash), []byte(m.Hash)) {
							http.Error(w, "wrong hash", http.StatusBadRequest)
							log.Printf("Wrong Hash of %s:counter:%d\n%s\n%s", m.ID, *m.Delta,
								m.Hash, hash)
							return
						}
					}
				} else {
					http.Error(w, "wrong type", http.StatusNotImplemented)
					log.Println("UpdatesJSONHandler:", m)
					return
				}

				if err := h.St.PutMetric(r.Context(), m.MType, m.ID, MetricValue); err != nil {
					//if err.Error() == "wrong type" {
					//	http.Error(w, err.Error(), http.StatusNotImplemented)
					//} else if err.Error() == "wrong value" {
					//	http.Error(w, err.Error(), http.StatusBadRequest)
					//} else {
					//	http.Error(w, "unknown error", http.StatusBadRequest)
					//}
					switch err.(type) {
					case *storage.TypeError:
						http.Error(w, err.Error(), http.StatusNotImplemented)
					default:
						http.Error(w, err.Error(), http.StatusBadRequest)
					}
					log.Println("UpdateJSONHandler:PutMetric ", err)
					return
				}
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

			v, err := h.St.GetMetric(r.Context(), m.MType, m.ID)
			if err != nil {
				//if err.Error() == "wrong type" {
				//	http.Error(w, err.Error(), http.StatusNotImplemented)
				//} else if strings.HasPrefix(err.Error(), "no such metric") {
				//	http.Error(w, err.Error(), http.StatusNotFound)
				//} else {
				//	http.Error(w, err.Error(), http.StatusBadRequest)
				//}
				switch err.(type) {
				case *storage.TypeError:
					http.Error(w, err.Error(), http.StatusNotImplemented)
				case *storage.MetricError:
					http.Error(w, err.Error(), http.StatusNotFound)
				default:
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
				if h.Key != "" {
					m.Hash = internal.Hash(fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value), h.Key)
				}
			} else if m.MType == "counter" {
				delta, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					log.Println(err)
				}
				m.Delta = &delta
				if h.Key != "" {
					m.Hash = internal.Hash(fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta), h.Key)
				}

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
