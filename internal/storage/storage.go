package storage

import "github.com/Feinot/metric-and-allert/internal/forms"

var Counter = make(map[string]int64)
var Gauge = make(map[string]float64)
var Storage forms.MemStorage
var M forms.Monitor
