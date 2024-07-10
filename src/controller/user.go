package controller

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	models "money-tracker-backend/src/model"
)

// TODO move into helpers later
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// isUniqueViolation checks if the error is due to a unique constraint violation
func isUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505" // PostgreSQL unique violation code
	}
	return false
}
func CreateUserHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user from the request payload
		var user models.User

		if err := c.Bind(&user); err != nil {
			fmt.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Hash the user's password
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			fmt.Println("Error hashing password:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		user.Password = hashedPassword

		// Log the received user
		fmt.Printf("Received user: %+v\n", user)

		// Query to search for user by username
		query := `
		INSERT INTO
			"users" (
				username,
				password,
				email,
				created_at,
				updated_at
			)
			VALUES ($1, $2, $3, NOW(), NOW())
			RETURNING
			id,
			username,
			email,
			created_at,
			updated_at
		`

		var newUser models.User
		err = db.QueryRow(query, user.Username, user.Password, user.Email).Scan(&newUser.ID, &newUser.Username, &newUser.Email, &newUser.CreatedAt, &newUser.UpdatedAt)
		if err != nil {
			fmt.Println("Error inserting user:", err)

			// Check if the error is due to unique constraint violation
			if isUniqueViolation(err) {
				response := WriteResponse{
					StatusCode: http.StatusConflict,
					Message:    "Username or email already exists",
					Data:       map[string]string{"error": err.Error()},
				}
				c.JSON(http.StatusConflict, response)
				return
			}

			response := WriteResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Could not create user",
				Data:       map[string]string{"error": err.Error()},
			}
			c.JSON(http.StatusConflict, response)
			return
		}
		// Log the received user
		fmt.Printf("Created user: %+v\n", newUser)
		response := WriteResponse{
			StatusCode: http.StatusOK,
			Message:    "Registration successful",
			Data:       newUser,
		}
		c.JSON(http.StatusOK, response)

	}
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

		var user models.User
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
			"message": "List Users",
			"data": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
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
		var users []models.User

		// Iterate over the result set
		for rows.Next() {
			var user models.User
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
		c.JSON(http.StatusOK, gin.H{
			"message": "List Users",
			"data":    users,
		})
	}
}
