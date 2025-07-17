package service

import (
	"fmt"
	"log"

	"example.com/mygamelist/repository"
)

type ListService struct {
	ListRepository *repository.ListRepository
}

func NewListService(repo *repository.ListRepository) *ListService {
	return &ListService{ListRepository: repo}
}

func (s *ListService) PostGame(gameId int, username, status string) error {
	err := s.ListRepository.InsertGame(gameId, username, status)
	if err != nil {
		log.Printf("failed to add game %s", err)
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (s *ListService) PutGame(gameId int, username, status string) error {
	err := s.ListRepository.UpdateGame(gameId, username, status)
	if err != nil {
		log.Printf("failed to update game %s", err)
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (s *ListService) GetGameList(username string) ([]repository.Game, error) {
	gamelist, err := s.ListRepository.FetchGames(username)
	if err != nil {
		log.Printf("failed to fetch gamelist from database %s", err)
		return nil, fmt.Errorf("%w", err)
	}
	return gamelist, nil
}
