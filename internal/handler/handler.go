package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Feinot/metric-and-allert/internal/forms"
	"github.com/Feinot/metric-and-allert/internal/storage"
)

type Metric forms.Metric

func HandleGuage(name string, value float64) *float64 {

	storage.ServerGauge[name] = value
	q := storage.ServerGauge[name]
	return &q

}
func HandleCaunter(name string, value int64) *int64 {

	if storage.ServerCounter[name] != 0 {
		storage.ServerCounter[name] += value

	} else {
		storage.ServerCounter[name] = value
	}
	q := storage.ServerCounter[name]
	return &q

}
func HandleUpdate(w http.ResponseWriter, r *http.Request) {

	var metrics forms.Metrics

	buf, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(buf, &metrics); err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch metrics.MType {
	case "gauge":
		//metrics.Value = new(float64)
		if *metrics.Value == 0 {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		metrics.Value = HandleGuage(metrics.ID, *metrics.Value)

		resp, err := json.Marshal(metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(resp)
	case "counter":
		if *metrics.Delta == 0 {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		metrics.Delta = HandleCaunter(metrics.ID, *metrics.Delta)

		resp, err := json.Marshal(metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(resp)
	}

}

func RequestUpdateHandle(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodPost:

		url := strings.Split(r.URL.Path, "/update/")
		url = strings.Split(url[1], "/")

		metricType := url[0]
		metricName := strings.TrimSpace(url[1])

		if metricName == "" {
			http.Error(w, "", http.StatusNotFound)

			return
		}

		switch metricType {
		case "gauge":

			if len(url) > 2 {
				url = strings.Split(url[2], "\n")
				value, err := strconv.ParseFloat(url[0], 64)

				if err != nil {
					fmt.Print(err)
					http.Error(w, "", http.StatusBadRequest)
					return
				}
				HandleGuage(metricName, value)
				http.Error(w, "", http.StatusOK)

			}

		case "counter":
			if len(url) > 2 {
				url = strings.Split(url[2], "\n")
				value, err := strconv.ParseInt(url[0], 10, 64)
				if err != nil {
					fmt.Print(err)
					http.Error(w, "", http.StatusBadRequest)
					return
				}
				HandleCaunter(metricName, value)
				w.Header().Set("Content-Type", "application/text")
				http.Error(w, "", http.StatusOK)
			}

		default:

			http.Error(w, "", http.StatusBadRequest)

			return
		}

	}
}
func HandleValue(w http.ResponseWriter, r *http.Request) {
	var metrics forms.Metrics
	var mt forms.VMetric
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	defer r.Body.Close()
	fmt.Println(buf.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &mt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metrics.ID = mt.ID
	metrics.MType = mt.MType
	fmt.Println(metrics.MType)
	switch metrics.MType {

	case "gauge":

		q := storage.ServerGauge[metrics.ID]
		metrics.Value = &q

		resp, err := json.Marshal(metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	case "counter":

		metrics.ID = mt.ID
		metrics.MType = mt.MType
		q := storage.ServerCounter[metrics.ID]
		metrics.Delta = &q

		resp, err := json.Marshal(metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	default:
		http.Error(w, "default", http.StatusBadRequest)
	}

}

func RequestValueHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		arr := make([]string, 3)
		url := strings.Split(r.URL.Path, "/value/")
		url = strings.Split(url[1], "/")

		copy(arr, url)
		metricType := arr[0]
		metricName := strings.TrimSpace(arr[1])

		if metricName == "" {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		switch metricType {
		case "gauge":

			if storage.ServerGauge[metricName] == 0 {
				http.Error(w, "", http.StatusNotFound)
				return
			}
			q := strconv.FormatFloat(storage.ServerGauge[metricName], 'f', -1, 64)
			w.Header().Set("Content-Type", "application/text")
			http.Error(w, q, http.StatusOK)
		case "counter":
			q := storage.ServerCounter[metricName]

			if q == 0 {
				http.Error(w, "", http.StatusNotFound)
				return

			}
			str := strconv.FormatInt(q, 10)
			w.Header().Set("Content-Type", "application/text")
			http.Error(w, str, http.StatusOK)

		default:
			http.Error(w, "", http.StatusNotFound)
			return
		}

	default:
		http.Error(w, "default", 400)
	}

}
func HomeHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmpl := template.Must(template.New("storage.Storage").Parse(`<div>
            <h1>Guage</h1>
			<p1>{{ .Guage}}</p>
			<h1>Counter</h1>
            <p>{{ .Counter}}</p>
        </div>`))
		tmpl.Execute(w, storage.Storage)
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
