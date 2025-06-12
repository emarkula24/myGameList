package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/mygamelist/handler"
	"example.com/mygamelist/repository"
	"example.com/mygamelist/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mock_api.go
type MockAPI struct {
	mock.Mock
}

func (m *MockAPI) SearchGames(query string) (*http.Response, error) {
	args := m.Called(query)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestSearchHandler(t *testing.T) {
	// Mock response body
	mockBody := `{"games": ["Game1", "Game2"]}`
	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(mockBody)),
	}

	mockAPI := new(MockAPI)
	mockAPI.On("SearchGames", "zelda").Return(response, nil)

	env := &utils.Env{
		API: mockAPI,
	}

	repo := &repository.Repository{} // stub if needed
	h := handler.NewHandler(env, repo)

	req := httptest.NewRequest("GET", "/games/search?query=zelda", nil)
	w := httptest.NewRecorder()

	h.Search(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	assert.JSONEq(t, mockBody, string(body))

	mockAPI.AssertExpectations(t)
}
