package service

import (
	"database/sql"
	"fmt"

	"example.com/mygamelist/repository"
	"example.com/mygamelist/utils"
	"github.com/golang-jwt/jwt/v5"
)

func RegisterUser(db *sql.DB, username, email, password string) (int64, error) {

	isUser := repository.SelectUserByEmail(db, username)

	if isUser {
		return 0, fmt.Errorf("user already exists")
	}

	hashedPassword, _ := utils.HashPassword(password)

	return repository.InsertUser(db, username, email, hashedPassword)

}

func LoginUser(db *sql.DB, username, password string) (string, error) {

	var secretKey = []byte("secret-key")
	hashedPassword, err := repository.PasswordByUsername(db, username, password)
	if err != nil {
		return "", err
	}
	if !utils.CheckPasswordHash(hashedPassword, password) {
		return "Failed to match password", fmt.Errorf("invalid username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"username": username,
	})

	tokenString, err := token.SignedString(secretKey)

	return tokenString, err

}
