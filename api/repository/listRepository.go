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

func (r *ListRepository) InsertGame(gameId int, username, status string) error {

	// ON DUPLICATE KEY is used to keep from inserting duplicates
	query := `INSERT INTO games (game_id) VALUES (?) ON DUPLICATE KEY UPDATE game_id = game_id`
	_, err := r.Db.Exec(query, gameId)
	if err != nil {
		return fmt.Errorf("failed to insert game into table games %w", err)
	}
	query = `INSERT INTO user_games (game_id, username, status) VALUES (?,?,?) ON DUPLICATE KEY UPDATE username = username`
	_, err = r.Db.Exec(query, gameId, username, status)
	if err != nil {
		return fmt.Errorf("failed to insert game into table user_games %w", err)
	}
	return nil
}

func (r *ListRepository) UpdateGame(gameId int, username, status string) error {
	query := `
			UPDATE user_games 
			SET status = ? 
			WHERE username = ? AND game_id = ?
			`
	_, err := r.Db.Exec(query, status, username, gameId)
	if err != nil {
		return fmt.Errorf("failed to update game (game_id=%d, user_id=%s): %w", gameId, username, err)
	}
	return nil
}

type Game struct {
	GameID int    `json:"id"`
	Status string `json:"status"`
}

func (r *ListRepository) FetchGames(username string) ([]Game, error) {
	query := `
			SELECT game_id, status
			FROM user_games
			WHERE username = ?
	`
	rows, err := r.Db.Query(query, username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gamelist (username=%s): %w", username, err)
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		var game Game
		if err := rows.Scan(&game.GameID, &game.Status); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		games = append(games, game)
	}

	return games, nil
}
