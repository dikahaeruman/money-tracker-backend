package repositories

import (
	"database/sql"
	"money-tracker-backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) (*models.User, error) {
	query := `
        INSERT INTO users (username, password, email, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id, username, email, created_at, updated_at
    `
	err := r.db.QueryRow(query, user.Username, user.Password, user.Email).
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

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	query := "SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = $1"
	var user models.User
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindPasswordByEmail(email string) (*models.User, error) {
	query := "SELECT password FROM users WHERE email = $1"
	var user models.User
	err := r.db.QueryRow(query, email).Scan(&user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
