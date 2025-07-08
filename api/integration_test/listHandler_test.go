package integration_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/mygamelist/handler"
	"example.com/mygamelist/repository"
	"example.com/mygamelist/routes"
	"example.com/mygamelist/service"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ListTestSuite struct {
	DB      *sql.DB
	Server  *httptest.Server
	Router  *mux.Router
	Handler *handler.ListHandler
}

func NewListTestSuite(db *sql.DB) *ListTestSuite {
	repo := repository.NewListRepository(db)
	service := service.NewListService(repo)
	handler := handler.NewListHandler(service)

	router := mux.NewRouter()
	router = routes.CreateListSubRouter(router, handler)
	server := httptest.NewServer(router)

	return &ListTestSuite{
		DB:      db,
		Server:  server,
		Handler: handler,
		Router:  router,
	}
}

var listTestSuite *ListTestSuite

func TestAddToList(t *testing.T) {

	body := `{
		"game_id":"34126",
		"status":"playing",
		"user_id": "1"
	}`

	r, err := http.Post(listTestSuite.Server.URL+"/list/add", "application/json", strings.NewReader(body))
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusOK, r.StatusCode)
}
