package models

import "time"

type User struct {
	ID        int       `db:"id" json:"id"`
	Username  string    `db:"username" json:"username" form:"username"`
	Password  string    `db:"password" json:"password" form:"password"`
	Email     string    `db:"email" json:"email" form:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at,omitempty" `
	UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
