package service

import "example.com/mygamelist/repository"

type ListService struct {
	ListRepository *repository.ListRepository
}

func NewListService(repo *repository.ListRepository) *ListService {
	return &ListService{ListRepository: repo}
}

func (s *ListService) AddGame(gameId, userId, status string) error {
	s.ListRepository.InsertGame(gameId, userId, status)
	return nil
}
