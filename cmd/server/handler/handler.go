package handler

import (
	"fmt"
	"github.com/Feinot/metric-and-allert/forms"
	"github.com/Feinot/metric-and-allert/storage"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type Metric forms.Metric

var m Metric

func HandleGuage(w http.ResponseWriter) {

	s := make(map[string]float64)
	s[m.MetricName] = m.Guage
	storage.Storage.Guage = s
	w.WriteHeader(200)

}
func HandleCaunter(w http.ResponseWriter) {
	s := make(map[string]int64)

	if storage.Storage.Counter[m.MetricName] != 0 {
		s[m.MetricName] = storage.Storage.Counter[m.MetricName] + m.Counter

	} else {
		s[m.MetricName] = m.Counter
	}
	storage.Storage.Counter = s
	w.WriteHeader(200)
}

func RequestUpdateHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var err error
		arr := make([]string, 3)
		url := strings.Split(r.URL.Path, "/update/")
		url = strings.Split(url[1], "/")

		copy(arr, url)

		m.MetricType = url[0]
		m.MetricName = url[1]
		if m.MetricName == "" {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		switch m.MetricType {
		case "gauge":
			if len(url) > 2 {
				url = strings.Split(url[2], "\n")
				m.Guage, err = strconv.ParseFloat(url[0], 64)
				fmt.Println(m.Guage)
				if err != nil {
					fmt.Print(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				HandleGuage(w)

			}

		case "counter":
			if len(url) > 2 {
				url = strings.Split(url[2], "\n")
				m.Counter, err = strconv.ParseInt(url[0], 10, 64)
				if err != nil {
					fmt.Print(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
			HandleCaunter(w)
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
		m.MetricType = arr[0]
		m.MetricName = arr[1]
		if m.MetricName == "" {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		switch m.MetricType {
		case "gauge":
			q := storage.Storage.Guage[m.MetricName]
			if q == 0 {
				http.Error(w, "", http.StatusNotFound)
				return
			}
			s := strconv.FormatFloat(q, 'f', 6, 64)
			http.Error(w, s[:len(s)-4], http.StatusOK)
		case "counter":
			q := storage.Storage.Counter[m.MetricName]
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
