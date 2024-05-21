package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
)

//export POSTGRESQL_URL='postgres://postgres@localhost:5432/money-tracker-db?sslmode=disable'
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
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

// LoginHandler handles the login request and generates JWT token
func LoginHandler(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// This is a sample check, replace it with actual user authentication logic
	if creds.Username != "user" || creds.Password != "password" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
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

func main() {
	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open the connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
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

	app.Use(sentrygin.New(sentrygin.Options{}))

	app.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world!")
	})

	app.POST("/login", LoginHandler)
	app.GET("/protected", ProtectedHandler)

	app.Run(":8080")
}
