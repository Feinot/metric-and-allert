package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/Feinot/metric-and-allert/internal/forms"
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
	storage.M.NumGoroutine.Value = float64(runtime.NumGoroutine())
	storage.M.Alloc.MName = "NumGoroutine"

	storage.M.Alloc.Value = float64(rtm.Alloc)
	storage.M.Alloc.MName = "Alloc"
	storage.M.BuckHashSys.Value = float64(rtm.BuckHashSys)
	storage.M.BuckHashSys.MName = "BuckHashSys"
	storage.M.GCCPUFraction.Value = rtm.GCCPUFraction
	storage.M.GCCPUFraction.MName = "GCCPUFraction"
	storage.M.TotalAlloc.Value = float64(rtm.TotalAlloc)
	storage.M.TotalAlloc.MName = "TotalAlloc"
	storage.M.Sys.Value = float64(rtm.Sys)
	storage.M.Sys.MName = "Sys"
	storage.M.Mallocs.Value = float64(rtm.Mallocs)
	storage.M.Mallocs.MName = "Mallocs"
	storage.M.Frees.Value = float64(rtm.Frees)
	storage.M.Frees.MName = "Frees"
	storage.M.GCSys.Value = float64(rtm.GCSys)
	storage.M.GCSys.MName = "GCSys"
	storage.M.HeapAlloc.Value = float64(rtm.HeapAlloc)
	storage.M.HeapAlloc.MName = "HeapAlloc"
	storage.M.HeapIdle.Value = float64(rtm.HeapIdle)
	storage.M.HeapIdle.MName = "HeapIdle"
	storage.M.HeapInuse.Value = float64(rtm.HeapInuse)
	storage.M.HeapInuse.MName = "HeapInuse"
	storage.M.HeapObjects.Value = float64(rtm.HeapObjects)
	storage.M.HeapObjects.MName = "HeapObjects"
	storage.M.HeapReleased.Value = float64(rtm.HeapReleased)
	storage.M.HeapReleased.MName = "HeapReleased"
	storage.M.HeapSys.Value = float64(rtm.HeapSys)
	storage.M.HeapSys.MName = "HeapSys"
	storage.M.Lookups.Value = float64(rtm.Lookups)
	storage.M.Lookups.MName = "Lookups"
	storage.M.MCacheInuse.Value = float64(rtm.MCacheInuse)
	storage.M.MCacheInuse.MName = "MCacheInuse"
	storage.M.MCacheSys.Value = float64(rtm.MCacheSys)
	storage.M.MCacheSys.MName = "MCacheSys"
	storage.M.MSpanInuse.Value = float64(rtm.MSpanInuse)
	storage.M.MSpanInuse.MName = "MSpanInuse"
	storage.M.MSpanSys.Value = float64(rtm.MSpanSys)
	storage.M.MSpanSys.MName = "MSpanSys"
	storage.M.NextGC.Value = float64(rtm.NextGC)
	storage.M.NextGC.MName = "NextGC"
	storage.M.NumForcedGC.Value = float64(rtm.NumForcedGC)
	storage.M.NumForcedGC.MName = "NumForcedGC"
	storage.M.OtherSys.Value = float64(rtm.OtherSys)
	storage.M.OtherSys.MName = "OtherSys"
	storage.M.PauseTotalNs.Value = float64(rtm.PauseTotalNs)
	storage.M.PauseTotalNs.MName = "PauseTotalNs"
	storage.M.StackInuse.Value = float64(rtm.StackInuse)
	storage.M.StackInuse.MName = "StackInuse"
	storage.M.StackSys.Value = float64(rtm.StackSys)
	storage.M.StackSys.MName = "StackSys"
	storage.M.NumGC.Value = float64(rtm.NumGC)
	storage.M.NumGC.MName = "NumGC"
	storage.M.PollCount.Value += 1
	storage.M.PollCount.MName = "PollCount"
	storage.M.RandomValue.Value = rand.Float64()
	storage.M.RandomValue.MName = "RandomValue"

}
func SandGaugeRequest(host string) error {

	resp, err := client.Post(fmt.Sprintf("%s%s%s%s%v", host, "/update/gauge/", storage.M.RandomValue.MName, "/", storage.M.RandomValue.Value), "text/plain", nil)

	if err != nil {
		return fmt.Errorf("cannot sand post request gauge: %w", err)
	}
	defer resp.Body.Close()
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
func SandJsonGaugeRequest(host string) error {
	var metrics forms.Metrics
	metrics.Value = &storage.M.RandomValue.Value
	metrics.ID = storage.M.RandomValue.MName
	metrics.MType = "gauge"
	sp, err := json.Marshal(metrics)
	q := bytes.NewReader(sp)
	if err != nil {
		fmt.Println("err")
		return err
	}
	fmt.Println(*metrics.Value)

	resp, err := client.Post(fmt.Sprintf("%s%s", host, "/update/"), "application/json", q)
	resp.Header.Set("Content-Type", "application/json")

	if err != nil {
		return fmt.Errorf("cannot sand post request gauge: %w", err)
	}
	defer resp.Body.Close()
	return nil

}
func SandJsonCounterRequest(host string) error {
	var metrics forms.Metrics
	metrics.Delta = &storage.M.PollCount.Value
	metrics.ID = storage.M.PollCount.MName
	metrics.MType = "counter"
	sp, err := json.Marshal(metrics)
	q := bytes.NewReader(sp)
	if err != nil {
		fmt.Println("err")
		return err
	}
	//fmt.Println(*metrics.Value)

	resp, err := client.Post(fmt.Sprintf("%s%s", host, "/update/"), "application/json", q)

	if err != nil {
		return fmt.Errorf("cannot sand post request counter: %w", err)
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
			err := SandJsonGaugeRequest(host)
			if err != nil {
				fmt.Print("cannot sand Gauge post request:", err)
			}
			err = SandJsonCounterRequest(host)
			if err != nil {
				fmt.Print("cannot sand Gauge post request: ", err)
			}

		}
	}
}
