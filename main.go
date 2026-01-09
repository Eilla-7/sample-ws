package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	InitDB()
	http.HandleFunc("/auth", AuthHandler)
	http.HandleFunc("/query", QueryHandler)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	if !VerifyUserInDB(input.Username, input.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// token, _ := GenerateToken(input.Username)

	token, err := GenerateToken(input.Username)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func QueryHandler(w http.ResponseWriter, r *http.Request) {

	// tokenStr := r.Header.Get("Authorization")

	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return
	}
	tokenStr := parts[1]
	claims, valid := ValidateToken(tokenStr)

	if !valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// json.NewEncoder(w).Encode(map[string]string{
	// 	"user": claims.Username,
	// 	"data": "data",
	// })

	data, err := GetUserData(claims.Username)
	if err != nil {
		http.Error(w, "No data found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": claims.Username,
		"data": data,
	})
}
