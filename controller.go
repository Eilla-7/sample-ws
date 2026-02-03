package main

import (
	"encoding/json"
	"net/http"
	"strings"
	// "golang.org/x/crypto/bcrypt"
)

// func getHashPassword(password string) (string, error) {
// 	bytePassword := []byte(password)
// 	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(hash), nil
// }

func NewUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if input.Username == "" || input.Password == "" {
		http.Error(w, "Userename and passord are requied", http.StatusBadRequest)
		return
	}
	if VerifyUserInDB(input.Username, input.Password) {
		var exists int
		db.QueryRow("SELECT 1 FROM user WHERE Username = ?", input.Username).Scan(&exists)
		http.Error(w, "User already exist!", http.StatusConflict)
		return
	}
	err = Rigester(input.Username, input.Password)
	if err != nil {
		http.Error(w, "Can't Create User", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully\n"))

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

	token, err := GenerateToken(input.Username)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func QueryHandler(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return
	}

	claims, valid := ValidateToken(parts[1])
	if !valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {

	case http.MethodGet:
		data, err := GetUserData(claims.Username)
		if err != nil {
			http.Error(w, "No data found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"user": claims.Username,
			"data": data,
		})

	case http.MethodPost:
		var input struct {
			Info string `json:"info"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil || input.Info == "" {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		err = InsertUserData(claims.Username, input.Info)
		if err != nil {
			http.Error(w, "Could not save data", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Data saved successfully"))

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
