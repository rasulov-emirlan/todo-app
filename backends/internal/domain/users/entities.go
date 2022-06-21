package users

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	accessLifeTime  = time.Minute * 10
	refreshLifeTime = time.Hour * 24 * 7

	RoleAdmin Role = iota
	RoleUser
)

type (
	Role uint

	User struct {
		ID           string `json:"id"`
		Username     string `json:"username"`
		Email        string `json:"email"`
		PasswordHash string `json:"-"`

		Role Role `json:"role"`

		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}
)

func comparePassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
