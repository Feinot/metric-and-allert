package main

import (
	"flag"
	"fmt"
	"github.com/Feinot/metric-and-allert/forms"
	"github.com/Feinot/metric-and-allert/storage"
	"log"
	"os"
	"strconv"

	"math/rand"
	"net/http"
	"runtime"

	"sync"
	"time"
)

type sum struct {
	mu sync.Mutex
	b  []byte
}
type Metric forms.Metric

var client = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {

		return nil
	},
}
var (
	p, r           int
	host           string
	Poll           int
	Met            Metric
	Sum            sum
	reportInterval = time.Duration(10) * time.Second

	interval = time.Duration(2) * time.Second
)

func GetMet() {

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
func MakeGURequest(host string) {

	body, err := client.Post(fmt.Sprintf("%s%s%s%s%v", host, "/update/gauge/", storage.M.RandomValue.MName, "/", storage.M.RandomValue.Value), "text/plain", nil)

	if err != nil {
		log.Fatal(err)
	}
	defer body.Body.Close()

}
func MakeCoRequest(host string) {

	body, err := client.Post(fmt.Sprintf("%s%s%s%s%v", host, "/update/counter/", storage.M.PollCount.MName, "/", storage.M.PollCount.Value), "text/plain", nil)

	if err != nil {
		log.Fatal(err)
	}
	defer body.Body.Close()

}
func GetConfigHost() {

	if os.Getenv("ADDRESS") != "" {
		host = os.Getenv("ADDRESS")
	}

}
func GetConfigReport() {

	intrv, err := strconv.Atoi(os.Getenv("REPORT_INTERVA"))
	if err == nil {
		p = intrv
	}

	flag.Parse()

}
func GetConfigPool() {

	intrv, err := strconv.Atoi(os.Getenv("POLL_INTERVAL"))
	if err == nil {
		r = intrv
	}

}
func main() {
	flag.IntVar(&r, "r", 2, "")
	flag.IntVar(&p, "p", 10, "")
	flag.StringVar(&host, "a", "localhost:8080", "")
	GetConfigHost()
	GetConfigReport()

	GetConfigPool()
	flag.Parse()

	reportInterval = time.Duration(p) * time.Second
	interval = time.Duration(r) * time.Second

	host = "http://" + host
	go Interval(host)
	select {}
}
func Interval(host string) {
	ticker := time.NewTicker(reportInterval)
	tick := time.NewTicker(interval)

	for {
		select {

		case <-tick.C:
			GetMet()
		case <-ticker.C:
			MakeGURequest(host)
			MakeCoRequest(host)

		}
	}
}
