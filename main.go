package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

//export POSTGRESQL_URL='postgres://postgres@localhost:5432/money-tracker-db?sslmode=disable'
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "money-tracker-db"
)

var jwtKey = []byte("my_secret_key")

// Credentials struct to handle username and password
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims struct to handle JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type User struct {
	ID       int
	Username string
	Password string
	Email    string
	// Other fields as needed
}

// LoginHandler handles the login request and generates JWT token
func LoginHandler(c *gin.Context) {
	var creds Credentials

	if err := c.BindJSON(&creds); err != nil {
		// Log the payload
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// This is a sample check, replace it with actual user authentication logic
	//TODO if username is found, and password hashed is correct, then return token
	// if creds.Username != "user" || creds.Password != "password" {
	// 	fmt.Printf("Received payload: %+v\n", creds)

	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	// 	return
	// }

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": tokenString})
}

// ProtectedHandler handles requests to the protected endpoint
func ProtectedHandler(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Welcome, %s!", claims.Username)})
}

// searchHandler handles the search request and returns user information
func SearchHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the username from the request payload
		var payload struct {
			Username string `json:"username"`
		}
		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		// Log the received payload
		fmt.Printf("Received payload: %+v\n", payload)

		// Query to search for user by username
		query := "SELECT id, username, password, email FROM users WHERE username = $1;"
		fmt.Println("Executing query:", query, "with username:", payload.Username)

		row := db.QueryRow(query, payload.Username)
		fmt.Printf("Query: %+v\n", row)

		var user User
		err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user"})
			}
			fmt.Printf("Error: %+v\n", err)

			return
		}

		// User found, return user details
		c.JSON(http.StatusOK, gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			// Do not return password in real applications
		})
	}
}

// searchHandler handles the search request and returns user information
func AllUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Query to select all users
		query := "SELECT id, username, password, email FROM users;"
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
			fmt.Printf("Error fetching users: %v\n", err)
			return
		}
		defer rows.Close()

		// Slice to hold users
		var users []User

		// Iterate over the result set
		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning user"})
				fmt.Printf("Error scanning user: %v\n", err)
				return
			}
			// Append user to slice
			users = append(users, user)
		}
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over users"})
			fmt.Printf("Error iterating over users: %v\n", err)
			return
		}

		// Return users in the response
		c.JSON(http.StatusOK, users)
	}
}

func main() {
	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	// Ensure the database connection is available
	if err := db.Ping(); err != nil {
		fmt.Println("Error pinging database:", err)
		return
	}

	fmt.Println("Successfully connected!")

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://aeb1b28f74b00ab978941319e3a675ee@o4507263852740608.ingest.us.sentry.io/4507263980994560",
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed : %v\n", err)
	}

	app := gin.Default()
	gin.SetMode(gin.DebugMode)
	app.Use(sentrygin.New(sentrygin.Options{}))

	app.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world!")
	})

	app.GET("/users", AllUser(db))
	app.POST("/search", SearchHandler(db))
	app.POST("/login", LoginHandler)
	app.GET("/protected", ProtectedHandler)

	app.Run(":8080")
}
