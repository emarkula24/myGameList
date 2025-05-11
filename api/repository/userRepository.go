package repository

import (
	"database/sql"
	"fmt"
	"time"
)

func SelectUserByUsername(db *sql.DB, username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS ( SELECT 1 FROM users WHERE username=?)`
	if err := db.QueryRow(query, username).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to retrieve user %s: %w", username, err)
	}
	return exists, nil
}

func InsertUser(db *sql.DB, username, email, hashedPassword string) (int64, error) {

	createdAt := time.Now()
	result, err := db.Exec(`INSERT INTO users (username, email, password, created_at) VALUES (?,?,?,?)`, username, email, hashedPassword, createdAt)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user into database: %w", err)
	}
	userId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert user ID: %w ", err)
	}
	return userId, nil
}

func PasswordByUsername(db *sql.DB, username, password string) (string, error) {

	var userPassword string
	query := `SELECT password FROM users WHERE username = ?`
	err := db.QueryRow(query, username).Scan(&userPassword)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve password: %w", err)
	}
	return userPassword, nil
}
