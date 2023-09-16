package form

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
