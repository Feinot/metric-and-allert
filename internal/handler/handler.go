package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/Feinot/metric-and-allert/internal/forms"
	"github.com/Feinot/metric-and-allert/internal/storage"
)

type Metric forms.Metric

var m Metric

func HandleGuage(name string, value float64) {

	storage.Gauge[name] = value

}
func HandleCaunter(name string, value int64) {

	if storage.Counter[name] != 0 {
		storage.Counter[name] += value

	} else {
		storage.Counter[name] = value
	}

}
func Handler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:

		url := strings.Split(r.URL.Path, "/update/")
		url = strings.Split(url[1], "/")

		metricType := url[0]
		metricName := strings.TrimSpace(url[1])

		if metricName == "" {
			http.Error(w, "met", http.StatusNotFound)

			return
		}

		switch metricType {
		case "gauge":
			if len(url) > 2 {
				url = strings.Split(url[2], "\n")
				value, err := strconv.ParseFloat(url[0], 64)

				if err != nil {
					fmt.Print(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				HandleGuage(metricName, value)
				w.WriteHeader(200)

			}

		case "counter":
			if len(url) > 2 {
				url = strings.Split(url[2], "\n")
				value, err := strconv.ParseInt(url[0], 10, 64)
				if err != nil {
					fmt.Print(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				HandleCaunter(metricName, value)
				w.WriteHeader(200)
			}

		default:

			w.WriteHeader(http.StatusBadRequest)

			return
		}
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

			if storage.Gauge[metricName] == 0 {
				http.Error(w, "", http.StatusNotFound)
				return
			}
			q := strconv.FormatFloat(storage.Gauge[metricName], 'f', 6, 64)

			http.Error(w, q[:len(q)-3], http.StatusOK)
		case "counter":
			q := storage.Counter[metricName]

			if q == 0 {
				http.Error(w, "", http.StatusNotFound)
				return

			}
			str := strconv.FormatInt(q, 10)

			http.Error(w, str, http.StatusOK)
		default:
			http.Error(w, "", http.StatusNotFound)
			return
		}

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
			http.Error(w, "met", http.StatusNotFound)

			return
		}

		switch metricType {
		case "gauge":
			if len(url) > 2 {
				url = strings.Split(url[2], "\n")
				value, err := strconv.ParseFloat(url[0], 64)

				if err != nil {
					fmt.Print(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				HandleGuage(metricName, value)
				w.WriteHeader(200)

			}

		case "counter":
			if len(url) > 2 {
				url = strings.Split(url[2], "\n")
				value, err := strconv.ParseInt(url[0], 10, 64)
				if err != nil {
					fmt.Print(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				HandleCaunter(metricName, value)
				w.WriteHeader(200)
			}

		default:

			w.WriteHeader(http.StatusBadRequest)

			return
		}
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

			if storage.Gauge[metricName] == 0 {
				http.Error(w, "", http.StatusNotFound)
				return
			}
			q := strconv.FormatFloat(storage.Gauge[metricName], 'f', 3, 64)

			http.Error(w, q, http.StatusOK)
		case "counter":
			q := storage.Counter[metricName]

			if q == 0 {
				http.Error(w, "", http.StatusNotFound)
				return

			}
			str := strconv.FormatInt(q, 10)

			http.Error(w, str, http.StatusOK)

		default:
			http.Error(w, "", http.StatusNotFound)
			return
		}
	case http.MethodPost:
		http.Error(w, "", http.StatusMethodNotAllowed)
		return

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

	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
