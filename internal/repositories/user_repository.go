package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"money-tracker-backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) (*models.User, error) {
	existingUserByEmail, err := r.FindUserByEmail(user.Email)
	if err != nil && err != sql.ErrNoRows {
		// If the error is due to a database error, handle it
		return nil, err
	}

	if existingUserByEmail != nil {
		return nil, fmt.Errorf("user with email %s already exists", user.Email)
	}

	existingUserByUsername, err := r.FindUserByUsername(user.Username)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingUserByUsername != nil {
		return nil, fmt.Errorf("user with username %s already exists", user.Username)
	}

	query := `
        INSERT INTO users (username, password, email, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id, username, email, created_at, updated_at
    `
	err = r.db.QueryRow(query, user.Username, user.Password, user.Email).
		Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	query := "SELECT id, username, email FROM users"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}(rows)

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	query := "SELECT id, username, email, created_at, updated_at FROM users WHERE email = $1"
	var user models.User
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindUserByUsername(username string) (*models.User, error) {
	query := "SELECT id, username, email, created_at, updated_at FROM users WHERE username = $1"
	var user models.User
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindPasswordByEmail(email string) (*models.User, error) {
	query := "SELECT id, email,password FROM users WHERE email = $1"
	var user models.User
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
