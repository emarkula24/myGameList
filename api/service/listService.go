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

func (s *ListService) PostGame(gameId, userId int, status string) error {
	err := s.ListRepository.InsertGame(gameId, userId, status)
	if err != nil {
		log.Printf("failed to add game %s", err)
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (s *ListService) PutGame(gameId, userId int, status string) error {
	err := s.ListRepository.UpdateGame(gameId, userId, status)
	if err != nil {
		log.Printf("failed to update game %s", err)
		return fmt.Errorf("%w", err)
	}
	return nil
}
