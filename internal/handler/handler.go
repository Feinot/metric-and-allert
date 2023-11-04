package handler

import (
	"bytes"
	"encoding/json"

	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Feinot/metric-and-allert/internal/forms"
	"github.com/Feinot/metric-and-allert/internal/logger"
	"github.com/Feinot/metric-and-allert/internal/storage"
	"github.com/go-chi/chi"
)

type Metric forms.Metric

func HandleGuage(name string, value float64) *float64 {

	storage.ServerGauge[name] = value
	q := storage.ServerGauge[name]
	if storage.Interval == 0 {
		storage.SaveMetrics(storage.File)
	}
	return &q

}
func HandleCaunter(name string, value int64) *int64 {

	if storage.ServerCounter[name] != 0 {
		storage.ServerCounter[name] += value

	} else {
		storage.ServerCounter[name] = value
	}
	if storage.Interval == 0 {
		storage.SaveMetrics(storage.File)
	}
	q := storage.ServerCounter[name]
	return &q

}
func HandleUpdate(w http.ResponseWriter, r *http.Request) {

	var metrics forms.Metrics

	buf, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		logger.LogError("Cannot io.ReadAll(r.Body)", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(buf, &metrics); err != nil {
		logger.LogError("Cannot Unmarshal: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch metrics.MType {
	case "gauge":

		if *metrics.Value == 0 {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		metrics.Value = HandleGuage(metrics.ID, *metrics.Value)

		resp, err := json.Marshal(metrics)
		if err != nil {
			logger.LogError("Cannot Marskal", err)
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
			logger.LogError("Cannot Marshal: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(resp)
	}

}

func RequestUpdateHandle(w http.ResponseWriter, r *http.Request) {

	metricType := chi.URLParam(r, "Mtype")
	metricName := strings.TrimSpace(chi.URLParam(r, "Mname"))
	metricValue := strings.TrimSpace(chi.URLParam(r, "Mvalue"))

	if metricName == "" {
		http.Error(w, "", http.StatusNotFound)

		return
	}

	switch metricType {
	case "gauge":

		if metricValue != "" {

			value, err := strconv.ParseFloat(metricValue, 64)

			if err != nil {
				logger.LogError("Cannot ParseFloat", err)
				http.Error(w, "", http.StatusBadRequest)
				return
			}
			HandleGuage(metricName, value)
			http.Error(w, "", http.StatusOK)
			return
		}
		http.Error(w, "", http.StatusBadRequest)
	case "counter":
		if metricValue != "" {

			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				logger.LogError("Cannot ParseInt: ", err)
				http.Error(w, "", http.StatusBadRequest)
				return
			}
			HandleCaunter(metricName, value)
			w.Header().Set("Content-Type", "application/text")
			http.Error(w, "", http.StatusOK)
		}
		http.Error(w, "", http.StatusBadRequest)
	default:

		http.Error(w, "", http.StatusBadRequest)

		return
	}

}

func HandleValue(w http.ResponseWriter, r *http.Request) {
	var metrics forms.Metrics
	var mt forms.VMetric
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	defer r.Body.Close()

	if err != nil {
		logger.LogError("Cannot read from Body: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &mt); err != nil {
		logger.LogError("Cannot unmarshal: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metrics.ID = mt.ID
	metrics.MType = mt.MType

	if metrics.ID == "" {
		http.Error(w, "", http.StatusNotFound)
		return

	}

	switch metrics.MType {

	case "gauge":
		if !storage.ContainsMetrics(metrics.ID, metrics.MType) {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		q := storage.ServerGauge[metrics.ID]
		metrics.Value = &q

		resp, err := json.Marshal(metrics)
		if err != nil {
			logger.LogError("Cannot Marshal: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	case "counter":

		q := storage.ServerCounter[metrics.ID]
		metrics.Delta = &q
		if !storage.ContainsMetrics(metrics.ID, metrics.MType) {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		resp, err := json.Marshal(metrics)
		if err != nil {
			logger.LogError("Cannot Marshal: ", err)
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

	metricType := chi.URLParam(r, "Mtype")
	metricName := strings.TrimSpace(chi.URLParam(r, "Mname"))

	if metricName == "" {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	if !storage.ContainsMetrics(metricName, metricType) {
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
		w.Header().Set("Content-Type", "html/text")
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
