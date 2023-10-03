package server

import (
	"fmt"
	"net/http"

	"github.com/Feinot/metric-and-allert/internal/config"
	"github.com/Feinot/metric-and-allert/internal/handler"
	"github.com/Feinot/metric-and-allert/internal/logger"
	"github.com/go-chi/chi"
)

func Run() {

	cfg := config.LoadServerConfig()
	r := chi.NewRouter()
	r.Use(logger.WithLogging)

	r.Post("/update/{type}/{name}/{value}", handler.RequestUpdateHandle)
	r.Post("/update/", handler.HabdleUpdate)

	r.Get("/value/{type}/{name}", handler.RequestValueHandle)
	r.Get("/value/", handler.RequestValueHandle)

	r.Get("/", (handler.HomeHandle))

	if err := http.ListenAndServe(cfg[1], r); err != nil {

		fmt.Println("Error: ", err)
	}
}

// logger.WithLogging
