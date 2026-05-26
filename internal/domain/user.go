package domain

import (
	"errors"
	"time"
)

type User struct {
	ID 	 string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	PasswordHash string `json:"password_hash"`
	CreatedAt time.Time `json:"created_at"`
}

var ErrUserNotFound = errors.New("user not found")

var ErrUserAlreadyExists = errors.New("user already exists")