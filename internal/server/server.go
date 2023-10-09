package server

import (
	"fmt"
	"net/http"

	"github.com/Feinot/metric-and-allert/internal/config"
	"github.com/Feinot/metric-and-allert/internal/gzipcompres"
	"github.com/Feinot/metric-and-allert/internal/handler"
	"github.com/Feinot/metric-and-allert/internal/logger"
	"github.com/Feinot/metric-and-allert/internal/storage"
	"github.com/go-chi/chi"
)

func Run() {

	cfg := config.LoadServerConfig()
	if cfg.Restore {
		storage.SelectMetric(cfg.File)
	}
	go storage.Run(cfg.File, cfg.Interval)
	r := chi.NewRouter()
	r.Use(logger.WithLogging)
	r.Use()

	r.Post("/update/{type}/{name}/{value}", handler.RequestUpdateHandle)
	r.Post("/update/", handler.HandleUpdate)

	r.Get("/value/{type}/{name}", handler.RequestValueHandle)
	r.Post("/value/", handler.HandleValue)

	r.Get("/", (handler.HomeHandle))

	if err := http.ListenAndServe(cfg.Host, gzipcompres.GzipMiddleware(r)); err != nil {
		storage.SaveMetrics(cfg.File)
		fmt.Println("Error: ", err)
	}
}

// logger.WithLogging
