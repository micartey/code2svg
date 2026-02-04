package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/svg", handleSVG)
	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
