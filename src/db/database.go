package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var DB *sql.DB // Exported global variable to hold the database connection

func ConnectDatabase() error {

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		os.Exit(1) // Exit the program with non-zero status code indicating failure
	}

	// Continue with your application logic
	fmt.Println("Environment variables loaded successfully!")

	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

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
