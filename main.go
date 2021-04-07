package main

import (
	"log"
	"net/http"

	"github.com/jaredpetersen/go-rest-example/internal/server"
)

func main() {
	srv := server.New()
	log.Fatal(http.ListenAndServe(":8080", srv))
}
