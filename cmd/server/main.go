package main

import (
	"fmt"
	"net/http"

	"problem_2/internal"
)

func main() {
	// Task 1
	http.HandleFunc("/search", internal.SearchHandler)

	// Task 2
	http.HandleFunc("/top-pairs", internal.TopPairsHandler)

	fmt.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server failed:", err)
	}
}