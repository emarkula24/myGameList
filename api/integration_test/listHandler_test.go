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
	client := &service.GiantBombClient{}
	repo := repository.NewListRepository(db)
	service := service.NewListService(repo, client)
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
func (s *ListTestSuite) GetServerURL() string {
	return s.Server.URL
}
func (s *ListTestSuite) GetClient() *http.Client {
	return s.Server.Client()
}

var listTestSuite *ListTestSuite

func TestAddToList(t *testing.T) {
	accessToken, _, _, err := RegisterAndLoginTestUser(listTestSuite, "listaddtest", "listadd@test.com", "1234567@M")
	require.NoError(t, err)

	body := `{
		"game_id":34126,
		"status":"playing",
		"username":"mies",
		"gamename":"metroid"
	}`
	client := listTestSuite.Server.Client()
	listAddRequest, err := http.NewRequest("POST", listTestSuite.Server.URL+"list/add", strings.NewReader(body))
	require.NoError(t, err)
	assert.NotNil(t, listAddRequest)
	listAddRequest.Header.Set("Authorization", accessToken)
	listAddRequest.Header.Set("Content-Type", "application/json")
	r, err := client.Do(listAddRequest)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, r.StatusCode)
}

func TestUpdateList(t *testing.T) {
	accessToken, _, _, err := RegisterAndLoginTestUser(listTestSuite, "listupdateest", "listupdate@test.com", "1234567@M")
	require.NoError(t, err)
	body := `{
		"game_id":34126,
		"status":"completed",
		"username":"mies"
	}`
	client := listTestSuite.Server.Client()
	updateReq, err := http.NewRequest("PUT", listTestSuite.Server.URL+"list/update", strings.NewReader(body))
	require.NoError(t, err)
	assert.NotNil(t, updateReq)
	updateReq.Header.Set("Authorization", accessToken)
	updateReq.Header.Set("Content-Type", "application/json")

	updateResp, err := client.Do(updateReq)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, updateResp.StatusCode)

}
