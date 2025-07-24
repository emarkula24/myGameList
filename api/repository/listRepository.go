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

func (r *ListRepository) InsertGame(gameId, status int, username, gamename string) error {

	// ON DUPLICATE KEY is used to keep from inserting duplicates
	query := `INSERT INTO games (game_id, gamename) VALUES (?,?) ON DUPLICATE KEY UPDATE game_id = game_id`
	_, err := r.Db.Exec(query, gameId, gamename)
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

func (r *ListRepository) UpdateGame(gameId, status int, username string) error {

	query := `
			UPDATE user_games 
			SET status = ? 
			WHERE username = ? AND game_id = ? AND status != ?
			`
	res, err := r.Db.Exec(query, status, username, gameId, status)
	if err != nil {
		return fmt.Errorf("failed to update game (game_id=%d, user_id=%s): %w", gameId, username, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated: likely no match for (username=%s, game_id=%d)", username, gameId)
	}
	return nil
}

type Game struct {
	GameID int    `json:"id"`
	Status string `json:"status"`
}

func (r *ListRepository) FetchGames(username string, page, limit int) ([]Game, error) {
	// Calculate the OFFSET
	offset := (page - 1) * limit
	query := `
			SELECT gm.game_id, ug.status
			FROM user_games ug
			JOIN games gm ON ug.game_id = gm.game_id
			WHERE ug.username = ?
			ORDER BY gm.gamename ASC
			LIMIT ?
			OFFSET ?
			`
	rows, err := r.Db.Query(query, username, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gamelist (username=%s): %w", username, err)
	}

	var games []Game
	for rows.Next() {
		var game Game
		if err := rows.Scan(&game.GameID, &game.Status); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		games = append(games, game)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return games, nil
}
func (r *ListRepository) FetchGame(username string, gameId int) *Game {
	query := `
			SELECT gm.game_id, ug.status
			FROM user_games ug
			JOIN games gm ON ug.game_id = gm.game_id
			WHERE ug.username = ? AND ug.game_id = ?
			ORDER BY gm.gamename ASC
	`
	var game Game
	err := r.Db.QueryRow(query, username, gameId).Scan(&game.GameID, &game.Status)
	if err != nil {
		return nil
	}
	return &game

}
