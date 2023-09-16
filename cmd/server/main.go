package main

import (
	"github.com/Feinot/metric-and-allert/cmd/server/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/update/", handler.RequestHandle)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
