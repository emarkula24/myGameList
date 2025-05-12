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
	expirationTime := time.Now().Add(15 * time.Minute).Unix()
	hashedPassword, err := repository.PasswordByUsername(db, username, password)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve users password: %w", err)
	}
	if !utils.CheckPasswordHash(hashedPassword, password) {
		return "", errorutils.ErrPasswordMatch
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      expirationTime,
	})

	tokenString, err := token.SignedString(secretKey)

	return tokenString, err

}

func StoreRefreshToken(db *sql.DB, username, refreshToken, jti string) error {
	userId, err := repository.SelectUserIdByUsername(db, username)
	if err != nil {
		return err
	}
	err = repository.InsertRefreshToken(db, userId, refreshToken, jti)
	if err != nil {
		return fmt.Errorf("failed to insert refreshtoken: %w", err)
	}

	return nil
}
