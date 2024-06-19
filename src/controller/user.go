package controller

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	models "money-tracker-backend/src/model"
)

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
