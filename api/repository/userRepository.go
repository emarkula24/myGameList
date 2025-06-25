package repository

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SelectUserByUsername(username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS ( SELECT 1 FROM users WHERE username=?)`
	if err := r.db.QueryRow(query, username).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to retrieve user %s: %w", username, err)
	}
	return exists, nil
}

func (r *Repository) SelectUserIdByUsername(username string) (int, error) {
	var userId int
	query := `SELECT user_id FROM users where username=?`
	if err := r.db.QueryRow(query, username).Scan(&userId); err != nil {
		return 0, fmt.Errorf("failed to retrieve userID from user %s: %w", username, err)
	}
	return userId, nil
}

func (r *Repository) InsertUser(username, email, hashedPassword string) (int64, error) {

	result, err := r.db.Exec(`INSERT INTO users (username, email, password) VALUES (?,?,?)`, username, email, hashedPassword)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user into database: %w", err)
	}
	userId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert user ID: %w ", err)
	}
	return userId, nil
}

func (r *Repository) PasswordByUsername(username, password string) (string, error) {

	var userPassword string
	query := `SELECT password FROM users WHERE username = ?`
	err := r.db.QueryRow(query, username).Scan(&userPassword)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve password: %w", err)
	}
	return userPassword, nil
}

func (r *Repository) InsertRefreshToken(userId int, refreshToken string, jti string) error {
	query := `
		INSERT INTO refreshtokens (user_id, refresh_token, jti)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
			refresh_token = VALUES(refresh_token),
			jti = VALUES(jti)
	`
	_, err := r.db.Exec(query, userId, refreshToken, jti)
	return err
}

func (r *Repository) RefreshTokenById(userId int) (string, string, error) {
	var token, jti string
	query := `SELECT refresh_token, jti from refreshtokens WHERE user_id = ?`
	err := r.db.QueryRow(query, userId).Scan(&token, &jti)
	if err != nil {
		return "", "", fmt.Errorf("failed to retrieve refresh token: %w", err)
	}

	return token, jti, nil
}
