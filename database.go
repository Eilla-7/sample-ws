package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID   int    `db:"UserID"`
	Username string `db:"Username"`
	Password string `db:"Password"`
}

type Query struct {
	QueryID int    `db:"QueryID"`
	UserID  int    `db:"UserID"`
	Info    string `db:"Info"`
}

var db *sql.DB

func InitDB() {
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = os.Getenv("DBNAME")

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

func Rigester(username, password string) error {
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}
	query := "INSERT INTO user (Username, Password) VALUES (?, ?)"
	_, err = db.Exec(query, username, string(hashed))
	if err != nil {
		return err
	}
	return nil
}

func VerifyUserInDB(username, password string) bool {
	var storedHash string
	query := "SELECT Password FROM user WHERE Username = ?"
	err := db.QueryRow(query, username).Scan(&storedHash)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	return err == nil
}

func InsertUserData(username, info string) error {
	var userID int

	err := db.QueryRow(
		"SELECT UserID FROM user WHERE Username = ?",
		username,
	).Scan(&userID)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"INSERT INTO query (UserId, Info) VALUES (?, ?)",
		userID,
		info,
	)
	return err
}

func GetUserData(username string) ([]string, error) {
	var userID int
	err := db.QueryRow("SELECT UserID FROM user WHERE Username = ?", username).Scan(&userID)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT Info FROM query WHERE UserId = ?", userID)
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
