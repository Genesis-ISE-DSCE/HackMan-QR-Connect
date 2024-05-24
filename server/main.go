package main

import (
	"fmt"
	"log"
	"net/http"
)

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func main() {

	http.HandleFunc("/", getHello)

	fmt.Println("Server is running on localhost:7500...")
    if err := http.ListenAndServe(":7500", nil); err!= nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}