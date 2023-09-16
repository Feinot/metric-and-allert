package forms

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
	Counter map[string][]int64
}
type GuageBody struct {
	Value float64
	MName string
}
type CounterBody struct {
	Value int64
	MName string
}
