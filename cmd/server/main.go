package main

import (
	"flag"
	"github.com/Feinot/metric-and-allert/cmd/server/handler"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	host string
)

func main() {

	Server()
}
func GetConfig() []string {

	if os.Getenv("ADDRESS") != "" {
		return strings.Split(os.Getenv("ADDRESS"), "localhost")
	}
	flag.StringVar(&host, "a", "localhost:8080", "")

	flag.Parse()
	return strings.Split(host, "localhost")
}
func Server() {

	q := GetConfig()
	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", handler.RequestUpdateHandle)

	r.Get("/value/{type}/{name}", handler.RequestValueHandle)

	r.Get("/", handler.HomeHandle)

	if err := http.ListenAndServe(q[1], r); err != nil {
		log.Fatal(err)
	}
}
