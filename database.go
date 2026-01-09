package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
	var err error
	dsn := "root:pssword@tcp(127.0.0.1:3306)/testdb"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
}

func VerifyUserInDB(username, password string) bool {
	var storedPass string
	query := "SELECT password FROM users WHERE username = ?"
	err := db.QueryRow(query, username).Scan(&storedPass)
	if err != nil {
		return false
	}
	return storedPass == password
}
