package service

import (
	"fmt"
	"log"
	"net/http"

	"example.com/mygamelist/interfaces"
	"example.com/mygamelist/repository"
)

type ListService struct {
	ListRepository *repository.ListRepository
	Cbc            interfaces.GiantBombClient
}

func NewListService(repo *repository.ListRepository, client interfaces.GiantBombClient) *ListService {
	return &ListService{ListRepository: repo, Cbc: client}
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

func (s *ListService) GetGameList(username string, page, limit int) (*http.Response, []repository.Game, error) {
	gamelist, err := s.ListRepository.FetchGames(username, page, limit)
	if err != nil {
		log.Printf("failed to fetch gamelist from database %s", err)
		return nil, nil, fmt.Errorf("%w", err)
	}
	fullGameList, err := s.Cbc.SearchGameList(gamelist, limit)
	if err != nil {
		log.Printf("failed to fetch gamelistdata from gamebomb %s", err)
		return nil, nil, fmt.Errorf("%w", err)
	}
	return fullGameList, gamelist, nil

}
