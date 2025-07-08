package integration_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"example.com/mygamelist/handler"
	"example.com/mygamelist/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
		mockError:      errors.New("giantbomb return != 200"),
		expectStatus:   http.StatusInternalServerError,
		expectContains: "failed to fetch gamedata",
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
		expectContains: "failed to fetch gamedata",
	},
	{
		name:       "API return 404",
		queryParam: "404",
		mockResponse: &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       nil,
		},
		mockError:      nil,
		expectStatus:   http.StatusInternalServerError,
		expectContains: "failed to fetch gamedata",
	},
}

// Unit tests for Gamebomb functionalities
func TestSearch(t *testing.T) {
	// External GameBomb API is mocked because without mocking it's not possible to produce a failure scenario.
	for _, tt := range searchTestCases {
		t.Run(tt.name, func(t *testing.T) {

			mockAPI := new(mocks.GameServiceMock)
			mockAPI.On("SearchGames", tt.queryParam).Return(tt.mockResponse, tt.mockError)

			h := handler.NewGameHandler(mockAPI)

			req, err := http.NewRequest(http.MethodGet, "/?query="+url.QueryEscape(tt.queryParam), nil)
			require.NoError(t, err)
			w := httptest.NewRecorder()

			h.Search(w, req)

			res := w.Result()
			assert.NotNil(t, res)
			err = res.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.expectStatus, res.StatusCode)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"), "Content type should be application json")
			bodyBytes, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			bodyStr := string(bodyBytes)
			assert.Contains(t, bodyStr, tt.expectContains)

			mockAPI.AssertExpectations(t)
		})
	}
}

type searchGameTestCase struct {
	name           string
	guid           string
	mockResponse   *http.Response
	mockError      error
	expectStatus   int
	expectContains string
}

var searchGameTestCases = []searchGameTestCase{
	{
		name: "Valid query returns OK",
		guid: "3030-6215",
		mockResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"status_code":1, "results": []}`)),
		},
		mockError:      nil,
		expectStatus:   http.StatusOK,
		expectContains: `"status_code":1`,
	},
	{
		name: "Empty query returns OK",
		guid: "",
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
		guid:           "errorcase",
		mockResponse:   nil,
		mockError:      errors.New("giantbomb return != 200"),
		expectStatus:   http.StatusInternalServerError,
		expectContains: "Failed to fetch gamedata",
	},
	{
		name: "API returns 500 on wrong-api-key error code",
		guid: "badstatus",
		mockResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"status_code":100}`)),
		},
		mockError:      nil,
		expectStatus:   http.StatusInternalServerError,
		expectContains: "Failed to fetch gamedata",
	},
	{
		name: "API return 404",
		guid: "404",
		mockResponse: &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       nil,
		},
		mockError:      nil,
		expectStatus:   http.StatusInternalServerError,
		expectContains: "Failed to fetch gamedata",
	},
}

func TestSearchGame(t *testing.T) {
	for _, tt := range searchGameTestCases {
		t.Run(tt.name, func(t *testing.T) {
			mockApi := new(mocks.GameServiceMock)
			mockApi.On("SearchGame", tt.guid).Return(tt.mockResponse, tt.mockError)

			h := handler.NewGameHandler(mockApi)

			req, err := http.NewRequest(http.MethodGet, "/games/game?guid="+url.QueryEscape(tt.guid), nil)
			require.NoError(t, err)
			w := httptest.NewRecorder()

			h.SearchGame(w, req)

			res := w.Result()
			assert.NotNil(t, res)

			require.NoError(t, err)
			assert.Equal(t, tt.expectStatus, res.StatusCode)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			bodyBytes, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			err = res.Body.Close()
			require.NoError(t, err)

			bodyStr := string(bodyBytes)
			assert.Contains(t, bodyStr, tt.expectContains)

			mockApi.AssertExpectations(t)
		})
	}

}
