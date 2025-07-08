package repository

import (
	"database/sql"
	"fmt"
)

type ListRepository struct {
	Db *sql.DB
}

func NewListRepository(Db *sql.DB) *ListRepository {
	return &ListRepository{Db: Db}
}

func (r *ListRepository) InsertGame(gameId, userId int, status string) error {

	query := `INSERT INTO games (game_id) VALUES (?) ON DUPLICATE KEY UPDATE game_id = game_id`
	_, err := r.Db.Exec(query, gameId)
	if err != nil {
		return fmt.Errorf("failed to insert game into table games %w", err)
	}
	query = `INSERT INTO user_games (game_id, user_id, status) VALUES (?,?,?) ON DUPLICATE KEY UPDATE user_id = user_id`
	_, err = r.Db.Exec(query, gameId, userId, status)
	if err != nil {
		return fmt.Errorf("failed to insert game into table user_games %w", err)
	}
	return nil
}
