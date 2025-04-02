package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Author: Aradhya Chakrabarti
// Roll No: 2205880
// Description: This program fetches numbers from external APIs and maintains a sliding window to compute the moving average.

// APIResponse struct to store the response from the external number API
type APIResponse struct {
	Numbers []int `json:"numbers"`
}

// SlidingWindow struct to store numbers and calculate moving average
type SlidingWindow struct {
	mu   sync.Mutex
	data []int
	size int
}

// Mapping number types to corresponding API URLs
var (
	numberAPIs = map[string]string{
		"p": "http://20.244.56.144/evaluation-service/primes",
		"f": "http://20.244.56.144/evaluation-service/fibo",
		"e": "http://20.244.56.144/evaluation-service/even",
		"r": "http://20.244.56.144/evaluation-service/rand",
	}
	windows = make(map[string]*SlidingWindow)
	token   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzQzNTk5NDkyLCJpYXQiOjE3NDM1OTkxOTIsImlzcyI6IkFmZm9yZG1lZCIsImp0aSI6ImQ5YmY0Mjg0LWM5OGYtNDc1Zi1hMmE4LTM1NDhmZWNiOGM4NiIsInN1YiI6ImFyYWRoeWEuY2hha3JhYmFydGlAZ21haWwuY29tIn0sImVtYWlsIjoiYXJhZGh5YS5jaGFrcmFiYXJ0aUBnbWFpbC5jb20iLCJuYW1lIjoiYXJhZGh5YSBjaGFrcmFiYXJ0aSIsInJvbGxObyI6IjIyMDU4ODAiLCJhY2Nlc3NDb2RlIjoibndwd3JaIiwiY2xpZW50SUQiOiJkOWJmNDI4NC1jOThmLTQ3NWYtYTJhOC0zNTQ4ZmVjYjhjODYiLCJjbGllbnRTZWNyZXQiOiJabW1GVnJzZkFSRmRYRFNuIn0.rWOzNdLLfkba_rdmamFkh8jAzDX2bNwNl5kDfc3CkGk"
)

// Function to fetch numbers from the external API
// Note that GET requests are not returning values:
/*
Sample curl command:
curl -X GET "http://20.244.56.144/evalua
tion-service/primes" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXV
CJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzQzNTk5NDkyLCJpYXQiOjE3NDM1OTkxOTIsImlzcyI6IkFm
Zm9yZG1lZCIsImp0aSI6ImQ5YmY0Mjg0LWM5OGYtNDc1Zi1hMmE4LTM1NDhmZWNiOGM4NiIsInN1YiI6
ImFyYWRoeWEuY2hha3JhYmFydGlAZ21haWwuY29tIn0sImVtYWlsIjoiYXJhZGh5YS5jaGFrcmFiYXJ0
aUBnbWFpbC5jb20iLCJuYW1lIjoiYXJhZGh5YSBjaGFrcmFiYXJ0aSIsInJvbGxObyI6IjIyMDU4ODAi
LCJhY2Nlc3NDb2RlIjoibndwd3JaIiwiY2xpZW50SUQiOiJkOWJmNDI4NC1jOThmLTQ3NWYtYTJhOC0z
NTQ4ZmVjYjhjODYiLCJjbGllbnRTZWNyZXQiOiJabW1GVnJzZkFSRmRYRFNuIn0.rWOzNdLLfkba_rdm
amFkh8jAzDX2bNwNl5kDfc3CkGk"
*/
// Output received:
// {"message":"Invalid authorization token"}
// So, the result is always going to show null

// For the purpose of testing, I have 2 test cases of my own.
func fetchNumbers(numberType string) ([]int, error) {
	url, exists := numberAPIs[numberType]
	if !exists {
		return nil, fmt.Errorf("invalid number type")
	}
	// Create HTTP request
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	// Execute HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Decode response
	var result APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Numbers, nil
}

// Adds numbers to the sliding window
func (w *SlidingWindow) Add(numbers []int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, num := range numbers {
		if len(w.data) >= w.size {
			w.data = w.data[1:]
		}
		w.data = append(w.data, num)
	}
}

// Computes the average of numbers in the sliding window
func (w *SlidingWindow) Average() float64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.data) == 0 {
		return 0
	}
	sum := 0
	for _, num := range w.data {
		sum += num
	}
	return float64(sum) / float64(len(w.data))
}

// Handler for processing number API requests
func numberHandler(w http.ResponseWriter, r *http.Request) {
	numberType := r.URL.Path[len("/numbers/"):] // Extract type from URL
	if _, exists := windows[numberType]; !exists {
		windows[numberType] = &SlidingWindow{size: 10}
	}
	numbers, err := fetchNumbers(numberType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	win := windows[numberType]
	win.Add(numbers)
	response := map[string]interface{}{
		"windowPrevState": win.data[:max(0, len(win.data)-len(numbers))],
		"windowCurrState": win.data,
		"numbers":         numbers,
		"avg":             win.Average(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func testNumbers() {
	testCases := map[string][]int{
		"p": {2, 3, 5, 7, 11},    // Prime numbers
		"f": {0, 1, 1, 2, 3, 5},  // Fibonacci numbers
		"e": {2, 4, 6, 8, 10},    // Even numbers
		"r": {7, 14, 21, 28, 35}, // Random numbers (mocked)
	}

	for key, numbers := range testCases {
		win := &SlidingWindow{size: 10}
		win.Add(numbers)
		response := map[string]interface{}{
			"numbers":         numbers,
			"windowPrevState": []int{}, // No previous state in test
			"windowCurrState": win.data,
			"avg":             win.Average(),
		}
		jsonResponse, _ := json.MarshalIndent(response, "", "  ")
		fmt.Printf("Test case %s:%s", key, jsonResponse)
	}
}

// Main function to start the HTTP server
func main() {
	go testNumbers()
	http.HandleFunc("/numbers/", numberHandler)
	fmt.Println("Server running on :9876")
	log.Fatal(http.ListenAndServe(":9876", nil))
}
