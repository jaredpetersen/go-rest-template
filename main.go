package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jaredpetersen/go-rest-example/internal/app"
)

func main() {
	a := app.New()
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      a,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
