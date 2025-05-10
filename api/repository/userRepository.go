package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func SelectUserByEmail(db *sql.DB, username string) bool {
	var exists bool
	query := `SELECT EXISTS ( SELECT 1 FROM users WHERE username=?)`
	if err := db.QueryRow(query, username).Scan(&exists); err != nil {
		log.Fatal(err)
	}
	return exists
}

func InsertUser(db *sql.DB, username, email, hashedPassword string) (int64, error) {

	createdAt := time.Now()
	result, err := db.Exec(`INSERT INTO users (username, email, password, created_at) VALUES (?,?,?,?)`, username, email, hashedPassword, createdAt)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	userId, _ := result.LastInsertId()
	return userId, nil
}

func PasswordByUsername(db *sql.DB, username, password string) (string, error) {

	var userPassword string
	query := `SELECT password FROM users WHERE username = ?`
	err := db.QueryRow(query, username).Scan(&userPassword)
	if err != nil {
		return "Failed to retrieve password", err
	}
	return userPassword, nil
}
