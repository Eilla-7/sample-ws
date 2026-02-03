package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error ...")
	}
	InitDB()
	http.HandleFunc("/signup", NewUserHandler)
	http.HandleFunc("/auth", AuthHandler)
	http.HandleFunc("/query", QueryHandler)

	port := os.Getenv("PORT")
	fmt.Println("Server listening on :", port)
	http.ListenAndServe(":"+port, nil)
}
