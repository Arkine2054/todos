package storer

import (
	"time"
)

type Users struct {
	ID           int        `db:"id"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	IsAdmin      bool       `db:"is_admin"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}

type ToDos struct {
	ID          int        `db:"id"`
	UserID      int        `db:"user_id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Completed   bool       `db:"completed"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

type Session struct {
	ID           string    `db:"id"`
	UserEmail    string    `db:"user_email"`
	RefreshToken string    `db:"refresh_token"`
	IsRevoked    bool      `db:"is_revoked"`
	CreatedAt    time.Time `db:"created_at"`
	ExpiresAt    time.Time `db:"expires_at"`
}

type List struct {
	Sort   string `db:"sort"`
	Order  string `db:"order"`
	Title  string `db:"title"`
	Limit  int    `db:"limit"`
	Offset int    `db:"offset"`
}
