package service

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"example.com/mygamelist/utils"
)

func AddUser(db *sql.DB, username, email, password string) (int64, error) {

	var exists bool
	query := `SELECT EXISTS ( SELECT 1 FROM users WHERE username=?)`
	if err := db.QueryRow(query, username).Scan(&exists); err != nil {
		log.Fatal(err)
	}
	if exists {
		return 0, fmt.Errorf("username already exists")
	}

	hashedPassword, _ := utils.HashPassword(password)
	createdAt := time.Now()
	result, err := db.Exec(`INSERT INTO users (username, email, password, created_at) VALUES (?,?,?,?)`, username, email, hashedPassword, createdAt)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	userId, _ := result.LastInsertId()
	return userId, nil
}
