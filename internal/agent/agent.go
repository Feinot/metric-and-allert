package agent

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/Feinot/metric-and-allert/internal/forms"
	"github.com/Feinot/metric-and-allert/internal/logger"
	"github.com/Feinot/metric-and-allert/internal/storage"
)

type Metric forms.Metric

var client = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {

		return nil
	},
}
var (
	p, r   int
	host   string
	Poll   int
	metric Metric
)

func GetMetric() {

	var rtm runtime.MemStats

	// Read full mem stats
	runtime.ReadMemStats(&rtm)

	// Number of goroutines

	storage.AgentGauge["NumGoroutine"] = float64(runtime.NumGoroutine()) //BuckHashSys
	storage.AgentGauge["Alloc"] = float64(rtm.Alloc)
	storage.AgentGauge["LastGC"] = float64(rtm.LastGC)
	storage.AgentGauge["BuckHashSys"] = float64(rtm.BuckHashSys)
	storage.AgentGauge["GCCPUFraction"] = rtm.GCCPUFraction
	storage.AgentGauge["TotalAlloc"] = float64(rtm.TotalAlloc)
	storage.AgentGauge["Sys"] = float64(rtm.Sys)
	storage.AgentGauge["Mallocs"] = float64(rtm.Mallocs)
	storage.AgentGauge["Frees"] = float64(rtm.Frees)
	storage.AgentGauge["GCSys"] = float64(rtm.GCSys)
	storage.AgentGauge["HeapAlloc"] = float64(rtm.HeapAlloc)
	storage.AgentGauge["HeapIdle"] = float64(rtm.HeapIdle)
	storage.AgentGauge["HeapInuse"] = float64(rtm.HeapInuse)
	storage.AgentGauge["HeapObjects"] = float64(rtm.HeapObjects)
	storage.AgentGauge["HeapReleased"] = float64(rtm.HeapReleased)
	storage.AgentGauge["HeapSys"] = float64(rtm.HeapSys)
	storage.AgentGauge["Lookups"] = float64(rtm.Lookups)
	storage.AgentGauge["MCacheInuse"] = float64(rtm.MCacheInuse)
	storage.AgentGauge["MCacheSys"] = float64(rtm.MCacheSys)
	storage.AgentGauge["MSpanInuse"] = float64(rtm.MSpanInuse)
	storage.AgentGauge["MSpanSys"] = float64(rtm.MSpanSys)
	storage.AgentGauge["NextGC"] = float64(rtm.NextGC)
	storage.AgentGauge["NumForcedGC"] = float64(rtm.NumForcedGC)
	storage.AgentGauge["OtherSys"] = float64(rtm.OtherSys)
	storage.AgentGauge["PauseTotalNs"] = float64(rtm.PauseTotalNs)
	storage.AgentGauge["StackInuse"] = float64(rtm.StackInuse)
	storage.AgentGauge["StackSys"] = float64(rtm.StackSys)
	storage.AgentGauge["NumGC"] = float64(rtm.NumGC)

	storage.AgentCounter["PollCount"] += 1
	storage.M.RandomValue.Value = rand.Float64()
	storage.AgentGauge["RandomValue"] = storage.M.RandomValue.Value

}
func SandGaugeRequest(host string) error {
	for key, value := range storage.AgentGauge {
		resp, err := client.Post(fmt.Sprintf("%s%s%s%s%v", host, "/update/gauge/", key, "/", value), "text/plain", nil)
		resp.Header.Set("Content-Encoding", "gzip")
		if err != nil {
			return fmt.Errorf("cannot sand post request gauge: %w", err)
		}
		defer resp.Body.Close()

	}
	return nil

}
func SandCounterRequest(host string) error {

	resp, err := client.Post(fmt.Sprintf("%s%s%s%s%v", host, "/update/counter/", storage.M.PollCount.MName, "/", storage.M.PollCount.Value), "text/plain", nil)

	if err != nil {
		return fmt.Errorf("cannot sand post request counter: %w", err)
	}
	defer resp.Body.Close()
	return nil

}
func SandJSONGaugeRequest(host string) error {
	for key, value := range storage.AgentGauge {
		metrics := forms.Metrics{
			Value: &value,
			ID:    key,
			MType: "gauge",
		}
		jMetrics, err := metrics.ToJason()
		if err != nil {
			return fmt.Errorf("cannot send request%v ", err)
		}
		SandJSON(jMetrics)

	}
	return nil
}
func SandJSONCounterRequest(host string) error {

	qas := storage.AgentCounter["PollCount"]
	metrics := forms.Metrics{
		Delta: &qas,
		ID:    "PollCount",
		MType: "counter",
	}
	jMetrics, err := metrics.ToJason()
	if err != nil {
		return fmt.Errorf("cannot send request%v ", err)
	}
	SandJSON(jMetrics)
	return nil

}
func SandJSON(sp []byte) error {
	var requestBody bytes.Buffer
	gz := gzip.NewWriter(&requestBody)

	if _, err := gz.Write([]byte(sp)); err != nil {
		logger.LogError("cannot create Writer: ", err)
		return fmt.Errorf("cannot marshal: %v   ", err)
	}

	if err := gz.Close(); err != nil {
		logger.LogError("cannot close NewWriter: ", err)
		return fmt.Errorf("cannot close NewWriter: %v   ", err)
	}

	url := fmt.Sprintf("%s%s", host, "/update/")

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		logger.LogError("cannot make POST request: ", err)
		return fmt.Errorf("cannot make POST request: %v   ", err)

	}
	req.Header.Set("Content-Encoding", "gzip")
	resp, err := client.Do(req)
	if err != nil {
		logger.LogError("cannot sand post request gauge: ", err)
		return fmt.Errorf("cannot sand post request gauge: %w  ", err)
	}
	defer resp.Body.Close()
	return nil

}
func Run(host string, reportInterval, interval time.Duration) {
	reportTicker := time.NewTicker(reportInterval)
	poolTicker := time.NewTicker(interval)

	for {
		select {

		case <-poolTicker.C:
			GetMetric()
		case <-reportTicker.C:

			if err := SandJSONGaugeRequest(host); err != nil {
				logger.LogError("cannot sand Gaug post request:", err)
			}

			if err := SandJSONCounterRequest(host); err != nil {
				logger.LogError("cannot sand Counter post request: ", err)
			}

		}
	}
}
