package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Repository struct {
	Db *sql.DB
}

func NewRepository(Db *sql.DB) *Repository {
	return &Repository{Db: Db}
}

func (r *Repository) SelectUserByUsername(username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS ( SELECT 1 FROM users WHERE username=?)`
	if err := r.Db.QueryRow(query, username).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to retrieve user %s: %w", username, err)
	}
	return exists, nil
}

func (r *Repository) SelectUserIdByUsername(username string) (int, error) {
	var userId int
	query := `SELECT user_id FROM users where username=?`
	if err := r.Db.QueryRow(query, username).Scan(&userId); err != nil {
		return 0, fmt.Errorf("failed to retrieve userID from user %s: %w", username, err)
	}
	return userId, nil
}

func (r *Repository) InsertUser(username, email, hashedPassword string) (int64, error) {

	result, err := r.Db.Exec(`INSERT INTO users (username, email, password, created_at) VALUES (?,?,?,?)`, username, email, hashedPassword, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to insert user into database: %w", err)
	}
	userId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert user ID: %w ", err)
	}
	return userId, nil
}

func (r *Repository) PasswordByUsername(username string) (string, error) {

	var userPassword string
	query := `SELECT password FROM users WHERE username = ?`
	err := r.Db.QueryRow(query, username).Scan(&userPassword)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve password: %w", err)
	}
	return userPassword, nil
}

func (r *Repository) InsertRefreshToken(userId int, refreshToken string, jti string) error {
	query := `
		INSERT INTO refreshtokens (user_id, refresh_token, jti, created_at)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			refresh_token = VALUES(refresh_token),
			jti = VALUES(jti)
	`
	_, err := r.Db.Exec(query, userId, refreshToken, jti, time.Now())
	return err
}

func (r *Repository) RefreshTokenById(userId int) (string, string, error) {
	var token, jti string
	query := `SELECT refresh_token, jti from refreshtokens WHERE user_id = ?`
	err := r.Db.QueryRow(query, userId).Scan(&token, &jti)
	if err != nil {
		return "", "", fmt.Errorf("failed to retrieve refresh token: %w", err)
	}

	return token, jti, nil
}
func (r *Repository) DeleteRefreshToken(userId int, jti string) error {
	query := `
			DELETE 
			FROM refreshtokens
			WHERE user_id = ? AND jti = ?
	`
	rows, err := r.Db.Exec(query, userId, jti)
	if err != nil {
		return err
	}
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		log.Printf("RowsAffected error: %v", err)
		return err
	}
	log.Println(rowsAffected)
	log.Println(rows)
	return err
}
