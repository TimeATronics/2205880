package main

import (
	"fmt"
	"log"
	"net/http"
)

// Author: Aradhya Chakrabarti
// Roll No: 2205880
// Description: This program fetches numbers from external APIs and maintains a sliding window to compute the moving average.

// APIResponse struct to store the response from the external number API
type APIResponse struct {
	Numbers []int `json:"numbers"`
}

func numberHandler(w http.ResponseWriter, r *http.Request) {}

// Main function to start the HTTP server
func main() {
	http.HandleFunc("/numbers/", numberHandler)
	fmt.Println("Server running on :9876")
	log.Fatal(http.ListenAndServe(":9876", nil))
}
