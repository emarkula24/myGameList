package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func AddUser(db *sql.DB, username, email, password string) (int64, error) {

	var exists bool
	query := `SELECT EXISTS ( SELECT 1 FROM users WHERE username=?)`
	if err := db.QueryRow(query, username).Scan(&exists); err != nil {
		log.Fatal(err)
	}
	if exists {
		return 0, fmt.Errorf("username already exists")
	}

	hashedPassword, _ := HashPassword(password)
	createdAt := time.Now()
	result, err := db.Exec(`INSERT INTO users (username, email, password, created_at) VALUES (?,?,?,?)`, username, email, hashedPassword, createdAt)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	userId, _ := result.LastInsertId()
	return userId, nil
}
