package db

import (
	"database/sql"
	"fmt"
	"os"
)

var DB *sql.DB // Exported global variable to hold the database connection

func ConnectDatabase() error {

	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	sslmode := os.Getenv("DATABASE_SSLMODE")
	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Connect to the database
	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}

	// Ping the database to verify connection
	if err := DB.Ping(); err != nil {
		fmt.Println("Error pinging database:", err)
		os.Exit(1)
	}
	fmt.Println("Successfully connected to the database!")
	return nil
}

func GetDB() *sql.DB {
	return DB
}
