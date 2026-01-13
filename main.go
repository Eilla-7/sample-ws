package main

import (
	"fmt"
	"net/http"
)

func main() {
	InitDB()
	http.HandleFunc("/auth", AuthHandler)
	http.HandleFunc("/query", QueryHandler)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
