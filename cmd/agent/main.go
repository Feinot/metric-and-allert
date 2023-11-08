package main

import (
	"github.com/Feinot/metric-and-allert/internal/agent"
	"github.com/Feinot/metric-and-allert/internal/config"

	"time"
)

func main() {

	cfg := config.LoadAgentConfig()

	reportInterval := time.Duration(cfg.ReportPool) * time.Second
	interval := time.Duration(cfg.Pool) * time.Second

	go agent.Run(cfg.Host, reportInterval, interval)
	select {}
}
