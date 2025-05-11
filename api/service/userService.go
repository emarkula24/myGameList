package service

import (
	"database/sql"
	"fmt"
	"time"

	"example.com/mygamelist/errorutils"
	"example.com/mygamelist/repository"
	"example.com/mygamelist/utils"
	"github.com/golang-jwt/jwt/v5"
)

func RegisterUser(db *sql.DB, username, email, password string) (int64, error) {

	isUser, err := repository.SelectUserByUsername(db, username)
	if err != nil {
		return 0, fmt.Errorf("failed to select user: %w", err)
	}

	if isUser {
		return 0, errorutils.ErrUserExists
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	userId, err := repository.InsertUser(db, username, email, hashedPassword)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	return userId, nil
}

func LoginUser(db *sql.DB, username, password string) (string, error) {

	var secretKey = []byte("secret-key")
	hashedPassword, err := repository.PasswordByUsername(db, username, password)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve users password: %w", err)
	}
	if !utils.CheckPasswordHash(hashedPassword, password) {
		return "", errorutils.ErrPasswordMatch
	}
	expirationTime := time.Now().Add(24 * time.Hour).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      expirationTime,
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign jwt token: %w", err)
	}

	return tokenString, err

}
