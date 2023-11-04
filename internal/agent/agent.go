package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
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

	storage.AgentCounter["PollCount"] = storage.M.PollCount.Value + 1
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
		var metrics forms.Metrics
		metrics.Value = &value
		metrics.ID = key
		metrics.MType = "gauge"
		var requestBody bytes.Buffer

		sp, err := json.Marshal(metrics)
		if err != nil {
			logger.LogError("Cannot unmarshal", err)
		}

		gz := gzip.NewWriter(&requestBody)
		_, err = gz.Write([]byte(sp))
		if err != nil {
			logger.LogError("cannot create writer", err)
		}
		err = gz.Close()
		if err != nil {
			logger.LogError("cannot close NewWriter", err)
			return err
		}

		fmt.Println(*metrics.Value)
		url := fmt.Sprintf("%s%s", host, "/update/")

		req, err := http.NewRequest("POST", url, &requestBody)
		if err != nil {
			logger.LogError("cannot make POST request ", err)

		}
		req.Header.Set("Content-Encoding", "gzip")
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("cannot sand post request gauge: %w%s   ", err, key)
		}
		defer resp.Body.Close()

	}
	return nil
}
func SandJSONCounterRequest(host string) error {
	var metrics forms.Metrics
	qas := storage.AgentCounter[storage.M.PollCount.MName]
	metrics.Delta = &qas
	metrics.ID = storage.M.PollCount.MName
	metrics.MType = "counter"
	sp, err := json.Marshal(metrics)
	var requestBody bytes.Buffer

	if err != nil {
		logger.LogError("cannot Marshal: ", err)
	}

	gz := gzip.NewWriter(&requestBody)
	_, err = gz.Write([]byte(sp))
	if err != nil {
		logger.LogError("cannot create Writer: ", err)
	}
	err = gz.Close()
	if err != nil {
		logger.LogError("cannot close NewWriter: ", err)
		return err
	}

	url := fmt.Sprintf("%s%s", host, "/update/")

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		logger.LogError("cannot make POST request: ", err)

	}
	req.Header.Set("Content-Encoding", "gzip")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot sand post request gauge: %w  ", err)
	}
	defer resp.Body.Close()
	return nil

}
func Run(host string, reportInterval, interval time.Duration) {
	ticker := time.NewTicker(reportInterval)
	tick := time.NewTicker(interval)

	for {
		select {

		case <-tick.C:
			GetMetric()
		case <-ticker.C:
			err := SandJSONGaugeRequest(host)
			if err != nil {
				logger.LogError("cannot sand Gauge post request:", err)
			}
			err = SandJSONCounterRequest(host)
			if err != nil {
				logger.LogError("cannot sand Counter post request: ", err)
			}

		}
	}
}
