package models

import (
	"errors"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("Username is required")
	}
	if u.Password == "" {
		return errors.New("Password is required")
	}
	if u.Email == "" {
		return errors.New("Email is required")
	}
	return nil
}
