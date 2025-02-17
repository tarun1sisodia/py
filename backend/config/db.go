package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// DB holds the database connection.
var DB *sql.DB

// ConnectDB connects to the MySQL database.
func ConnectDB() {
	var err error
	DB, err = sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Could not ping the database: %v", err)
	}

	log.Println("Connected to the database")
}

// CloseDB closes the database connection.
func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Error closing the database connection: %v", err)
		}
		log.Println("Database connection closed")
	}
}
