package main

import (
	"fmt"
	"net/http"
	"handler"
)

func main() {
    fmt.Printf("Starting...\n")

	http.HandleFunc("/location/", handler.LocationElement)
	http.HandleFunc("/location", handler.LocationIndex)
    http.ListenAndServe(":8080", nil)
}