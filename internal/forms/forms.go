package forms

import (
	"encoding/json"
	"fmt"

	"github.com/Feinot/metric-and-allert/internal/logger"
)

type Monitor struct {
	Alloc,
	TotalAlloc,
	Sys,
	Mallocs,
	Frees,
	LiveObjects,
	PauseTotalNs,
	BuckHashSys,
	GCSys,
	HeapAlloc,
	HeapIdle,
	HeapInuse,
	HeapObjects,
	HeapReleased,
	HeapSys,
	LastGC,
	Lookups,
	MCacheInuse,
	MCacheSys,
	MSpanInuse,
	MSpanSys,
	NextGC,
	OtherSys,
	StackInuse,
	StackSys,
	GCCPUFraction,
	RandomValue,
	NumForcedGC,
	NumGoroutine,
	NumGC GuageBody

	PollCount CounterBody
}
type Metric struct {
	MetricType string `json:"MetricType"`
	MetricName string `json:"MetricName"`
	Guage      float64
	Counter    int64
}
type MemStorage struct {
	Guage   map[string]float64
	Counter map[string]int64
}
type GuageBody struct {
	Value float64
	MName string
}
type CounterBody struct {
	Value int64
	MName string
}
type LoggerBody struct {
	URL          string
	Method       string
	Duration     int64
	StatusCode   int
	SizeResponse int
	TypeLog      string
}
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
type VMetric struct {
	ID    string `json:"id"` // имя метрики
	MType string `json:"type"`
}

func (m *Metrics) ToJason() ([]byte, error) {
	sp, err := json.Marshal(m)

	if err != nil {
		logger.LogError("cannot Marshal: ", err)
		return nil, fmt.Errorf("cannot marshal: %v   ", err)
	}
	return sp, nil
}
