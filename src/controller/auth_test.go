package controller

import (
	"bytes"
	"encoding/json"
	middleware "money-tracker-backend/src/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock the database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock user data
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	rows := sqlmock.NewRows([]string{"id", "username", "password", "email"}).
		AddRow(1, "testuser", string(hashedPassword), "test@example.com")

	mock.ExpectQuery("SELECT id, username, password, email FROM users WHERE email = \\$1").
		WithArgs("test@example.com").
		WillReturnRows(rows)

	// Mock JWT key
	jwtKey := []byte("my_secret_key")

	// Create a new HTTP request
	creds := Credentials{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonValue, _ := json.Marshal(creds)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	w := httptest.NewRecorder()

	// Create a new Gin router
	r := gin.Default()
	r.POST("/login", LoginHandler(db, jwtKey))

	// Serve the HTTP request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Login successful", response["message"])

	// Check the token in the response
	tokenString := response["token"]
	claims := &middleware.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.WithinDuration(t, time.Now().Add(5*time.Minute), time.Unix(claims.ExpiresAt, 0), time.Second)

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
