package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"pw_hash"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type CreateUserRequest struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

func NewUserFromRequest(r CreateUserRequest) (*User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	createdAt := time.Now().String()
	return &User{
		ID:           uuid.New().String(),
		Name:         r.Name,
		Email:        r.Email,
		PasswordHash: string(hashed),
		CreatedAt:    createdAt,
		UpdatedAt:    createdAt,
	}, nil
}
