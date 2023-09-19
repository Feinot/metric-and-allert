package server

import (
	"log"
	"net/http"

	"github.com/Feinot/metric-and-allert/internal/config"
	"github.com/Feinot/metric-and-allert/internal/handler"
	"github.com/go-chi/chi"
)

func Run() {

	q := config.LoadServerConfig()
	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", handler.RequestUpdateHandle)

	r.Get("/value/{type}/{name}", handler.RequestValueHandle)

	r.Get("/", handler.HomeHandle)

	if err := http.ListenAndServe(q[1], r); err != nil {
		log.Fatal(err)
	}
}
