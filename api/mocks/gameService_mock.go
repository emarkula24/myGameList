package mocks

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// MockAPI mocks the SearchGames method
type GameServiceMock struct {
	mock.Mock
}

func (m *GameServiceMock) SearchGames(query string) (*http.Response, error) {

	args := m.Called(query)
	resp, _ := args.Get(0).(*http.Response)
	return resp, args.Error(1)
}

func (m *GameServiceMock) SearchGame(guid string) (*http.Response, error) {
	args := m.Called(guid)
	return args.Get(0).(*http.Response), args.Error(1)
}
