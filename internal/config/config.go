package config

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

type AgentConfig struct {
	Pool       int
	ReportPool int
	Host       string
}
type ServerConfig struct {
	Host     string
	Interval int
	File     string
	Restore  bool
}

func LoadAgentConfig() AgentConfig {
	var config AgentConfig
	flag.IntVar(&config.Pool, "r", 2, "")
	flag.IntVar(&config.ReportPool, "p", 10, "")
	flag.StringVar(&config.Host, "a", "localhost:8080", "")
	flag.Parse()

	if os.Getenv("ADDRESS") != "" {
		config.Host = os.Getenv("ADDRESS")
	}
	r, err := strconv.Atoi(os.Getenv("POLL_INTERVAL"))
	if err == nil {
		config.Pool = r
	}
	p, err := strconv.Atoi(os.Getenv("REPORT_INTERVA"))
	if err == nil {
		config.ReportPool = p
	}
	config.Host = "http://" + config.Host
	return config

}
func LoadServerConfig() *ServerConfig {
	var config ServerConfig
	flag.IntVar(&config.Interval, "i", 300, "")
	flag.StringVar(&config.Host, "a", "localhost:8080", "")
	flag.StringVar(&config.File, "f", "/tmp/metrics-db.json", "")
	flag.BoolVar(&config.Restore, "r", true, "")
	flag.Parse()
	if os.Getenv("ADDRESS") != "" {
		config.Host = strings.Split(os.Getenv("ADDRESS"), "localhost")[1]
	}
	if os.Getenv("FILE_STORAGE_PATH") != "" {
		config.File = os.Getenv("FILE_STORAGE_PATH")
	}
	if os.Getenv("RESTORE") != "" {
		config.Restore, _ = strconv.ParseBool(os.Getenv("RESTORE"))
	}
	if os.Getenv("STORE_INTERVAL") != "" {
		config.Interval, _ = strconv.Atoi(os.Getenv("STORE_INTERVAL"))
	}

	return &config

}
