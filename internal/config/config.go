package config

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Pool       int
	ReportPool int
	Host       string
}

func LoadAgentConfig() Config {
	var config Config
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
func LoadServerConfig() []string {
	host := ""
	if os.Getenv("ADDRESS") != "" {
		return strings.Split(os.Getenv("ADDRESS"), "localhost")
	}
	flag.StringVar(&host, "a", "localhost:8080", "")

	flag.Parse()
	return strings.Split(host, "localhost")

}
