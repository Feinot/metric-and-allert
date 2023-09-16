package handler

import (
	"github.com/Feinot/metric-and-allert/cmd/server/form"
	"net/http"
	"strconv"
	"strings"
)

var storage form.MemStorage

type Metric form.Metric

var m Metric

func HandleGuage(w http.ResponseWriter) {

	s := make(map[string]float64)
	s[m.MetricName] = m.Guage
	storage.Guage = s
	http.Error(w, "", 200)

}
func HandleCaunter(w http.ResponseWriter) {

	s := make(map[string][]int64)
	s[m.MetricName] = append(storage.Counter[m.MetricName], m.Counter)
	storage.Counter = s
	http.Error(w, "", 200)
}

func RequestHandle(w http.ResponseWriter, r *http.Request) {
	var err error
	arr := make([]string, 3)
	url := strings.Split(r.URL.Path, "/update/")
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
		m.Guage, err = strconv.ParseFloat(arr[2], 64)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		HandleGuage(w)
	case "counter":

		m.Counter, err = strconv.ParseInt(arr[2], 10, 64)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		HandleCaunter(w)
	default:

		http.Error(w, "", http.StatusBadRequest)
		return
	}
}
