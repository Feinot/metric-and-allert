package main

import (
	"github.com/Feinot/metric-and-allert/cmd/server/handler"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {

	Server()
}
func Server() {
	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", handler.RequestUpdateHandle)

	r.Get("/value/{type}/{name}", handler.RequestValueHandle)

	r.Get("/", handler.HomeHandle)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
