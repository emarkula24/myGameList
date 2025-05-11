package service

import (
	"database/sql"
	"fmt"

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
		return "", err
	}
	if !utils.CheckPasswordHash(hashedPassword, password) {
		return "", errorutils.ErrPasswordMatch
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"username": username,
	})

	tokenString, err := token.SignedString(secretKey)

	return tokenString, err

}
