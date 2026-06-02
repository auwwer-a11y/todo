package domain

import (
	"errors"
	"time"
)

type User struct {
	ID 	 string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

var ErrUserNotFound = errors.New("user not found")

var ErrUserAlreadyExists = errors.New("user already exists")

var ErrInvalidCredentials = errors.New("invalid credentials")