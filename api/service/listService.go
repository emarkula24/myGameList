package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"example.com/mygamelist/repository"
)

// ListService defines a list service controller.
type ListService struct {
	ListRepository *repository.ListRepository
	Cbc            GameListFetcher
}

type GameListFetcher interface {
	SearchGameList(ctx context.Context, games []repository.Game, limit int) (*http.Response, error)
}

// NewListService creates a list service controller.
func NewListService(repo *repository.ListRepository, client GameListFetcher) *ListService {
	return &ListService{ListRepository: repo, Cbc: client}
}

// PostGame adds given game to list.
func (s *ListService) PostGame(gameId, status int, username, gamename string) error {
	return s.ListRepository.InsertGame(gameId, status, username, gamename)
}

// PutGame writes a game for a given list.
func (s *ListService) PutGame(gameId, status int, username string) error {
	return s.ListRepository.UpdateGame(gameId, status, username)
}

// GetGameList returns list for a given user.
func (s *ListService) GetGameList(ctx context.Context, username string, page, limit int) (*http.Response, []repository.Game, error) {
	gamelist, err := s.ListRepository.FetchGames(username, page, limit)
	if err != nil {
		log.Printf("failed to fetch gamelist from database %s", err)
		return nil, nil, fmt.Errorf("%w", err)
	}
	fullGameList, err := s.Cbc.SearchGameList(ctx, gamelist, limit)
	if err != nil {
		log.Printf("failed to fetch gamelistdata from gamebomb %s", err)
		return nil, nil, fmt.Errorf("%w", err)
	}

	return fullGameList, gamelist, nil

}

// GetGameFromList returns a game for a given list.
func (s *ListService) GetGameFromList(username string, gameId int) *repository.Game {
	return s.ListRepository.FetchGame(username, gameId)
}

// DeleteGame deletes a game from a given list.
func (s *ListService) DeleteGame(username string, gameId int) error {
	return s.ListRepository.RemoveGame(username, gameId)
}
