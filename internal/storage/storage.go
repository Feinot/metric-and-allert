package storage

import "github.com/Feinot/metric-and-allert/internal/forms"

var ServerCounter = make(map[string]int64)
var ServerGauge = make(map[string]float64)
var AgentCounter = make(map[string]int64)
var AgentGauge = make(map[string]float64)
var Storage forms.MemStorage
var M forms.Monitor
