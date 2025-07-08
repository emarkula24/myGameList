package service

import (
	"log"

	"example.com/mygamelist/repository"
)

type ListService struct {
	ListRepository *repository.ListRepository
}

func NewListService(repo *repository.ListRepository) *ListService {
	return &ListService{ListRepository: repo}
}

func (s *ListService) AddGame(gameId, userId int, status string) error {
	err := s.ListRepository.InsertGame(gameId, userId, status)
	if err != nil {
		log.Printf("failed to add game %s", err)
		return err
	}
	return nil
}
