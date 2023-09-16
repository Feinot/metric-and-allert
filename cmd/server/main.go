package main

import (
	"flag"
	"github.com/Feinot/metric-and-allert/cmd/server/handler"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strings"
)

var (
	host string
)

func main() {

	Server()
}
func Server() {
	flag.StringVar(&host, "a", "localhost:8080", "")

	flag.Parse()
	q := strings.Split(host, "localhost")

	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", handler.RequestUpdateHandle)

	r.Get("/value/{type}/{name}", handler.RequestValueHandle)

	r.Get("/", handler.HomeHandle)

	if err := http.ListenAndServe(q[1], r); err != nil {
		log.Fatal(err)
	}
}
