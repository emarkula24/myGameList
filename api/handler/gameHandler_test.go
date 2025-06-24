package handler_test

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"example.com/mygamelist/handler"
	"example.com/mygamelist/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAPI mocks the SearchGames method
type MockAPI struct {
	mock.Mock
}

type searchTestCase struct {
	name           string
	queryParam     string
	mockResponse   *http.Response
	mockError      error
	expectStatus   int
	expectContains string
}

var searchTestCases = []searchTestCase{
	{
		name:       "Valid query returns OK",
		queryParam: "mario",
		mockResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"status_code":1, "results": []}`)),
		},
		mockError:      nil,
		expectStatus:   http.StatusOK,
		expectContains: `"status_code":1`,
	},
	{
		name:       "Empty query returns OK",
		queryParam: "",
		mockResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"status_code":1, "results": []}`)),
		},
		mockError:      nil,
		expectStatus:   http.StatusOK,
		expectContains: `"status_code":1`,
	},
	{
		name:           "API returns 500 on error",
		queryParam:     "errorcase",
		mockResponse:   nil,
		mockError:      errors.New("API failure"),
		expectStatus:   http.StatusInternalServerError,
		expectContains: "Failed to fetch gamedata",
	},
	{
		name:       "API returns 500 on wrong-api-key error code",
		queryParam: "badstatus",
		mockResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"status_code":100}`)),
		},
		mockError:      nil,
		expectStatus:   http.StatusInternalServerError,
		expectContains: "Failed to fetch gamedata",
	},
}

func (m *MockAPI) SearchGames(query string) (*http.Response, error) {

	args := m.Called(query)
	resp, _ := args.Get(0).(*http.Response)
	return resp, args.Error(1)
}

// Unit tests for Gamebomb functionalities, logs provided but the log package are disables for the duration of test.
func TestSearchHandler(t *testing.T) {
	// External GameBomb API is mocked because without mocking it's not possible to produce a failure scenario.
	for _, tt := range searchTestCases {
		t.Run(tt.name, func(t *testing.T) {

			log.SetOutput(io.Discard)

			mockAPI := new(MockAPI)
			mockAPI.On("SearchGames", tt.queryParam).Return(tt.mockResponse, tt.mockError)

			env := &utils.Env{
				API: mockAPI,
			}
			h := handler.NewGameHandler(env)

			req, err := http.NewRequest(http.MethodGet, "/search?query="+url.QueryEscape(tt.queryParam), nil)
			assert.NoError(t, err, "creating request should not fail")
			w := httptest.NewRecorder()

			h.Search(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectStatus, res.StatusCode, "Test %q: status code mismatch:", tt.name)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"),
				"Content type should be application json")
			bodyBytes, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}

			bodyStr := string(bodyBytes)
			assert.Contains(t, bodyStr, tt.expectContains, "Test %q: response body missing expected substring:", tt.name)

			mockAPI.AssertExpectations(t)
		})
	}
}
