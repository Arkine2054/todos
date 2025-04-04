package handler

import (
	"time"
)

type UsersReq struct {
	Email        string     `json:"email" validate:"required,email"`
	PasswordHash string     `json:"password" validate:"required, min=0"`
	IsAdmin      bool       `json:"is_admin"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
type UsersRes struct {
	ID        int        `json:"id"`
	Email     string     `json:"email"`
	IsAdmin   bool       `json:"is_admin"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type ListUserRes struct {
	Users []UsersRes `json:"users"`
}

type LoginUserReq struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
}

type LoginUserRes struct {
	SessionID             string    `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  UsersRes  `json:"user"`
}

type RenewAccessTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}
type RenewAccessTokenRes struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

type ToDosReq struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type ToDosRes struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
type ListToDosRes struct {
	UserID int        `json:"user_id"`
	ToDos  []ToDosRes `json:"todos"`
}

type List struct {
	ID     int    `json:"id"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Title  string `json:"title"`
}
