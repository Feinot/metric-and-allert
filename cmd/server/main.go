package main

import (
	"github.com/Feinot/metric-and-allert/internal/logger"
	"github.com/Feinot/metric-and-allert/internal/server"
)

func main() {
	logger.Init()
	server.Run()
}
