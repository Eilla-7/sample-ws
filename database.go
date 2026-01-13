package main

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
	var err error
	dsn := "root:password@tcp(127.0.0.1:3306)/testdb"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
}

func VerifyUserInDB(username, password string) bool {
	var storedHash string
	query := "SELECT password FROM users WHERE username = ?"
	err := db.QueryRow(query, username).Scan(&storedHash)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	return err == nil
}


func GetUserData(username string) ([]string, error) {
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT info FROM data WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var info string
		rows.Scan(&info)
		results = append(results, info)
	}
	return results, nil
}
