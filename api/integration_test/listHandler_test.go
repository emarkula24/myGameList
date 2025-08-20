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
	"example.com/mygamelist/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ListTestSuite struct {
	DB          *sql.DB
	Server      *httptest.Server
	Router      *mux.Router
	ListHandler *handler.ListHandler
	UserHandler *handler.UserHandler
}

func NewListTestSuite(db *sql.DB) *ListTestSuite {
	client := &service.GiantBombClient{}
	repo := repository.NewListRepository(db)
	listService := service.NewListService(repo, client)
	listHandler := handler.NewListHandler(listService)

	userRepo := repository.NewRepository(db)
	auth := &utils.AuthService{}
	userService := service.NewUserService(userRepo, auth)
	userHandler := handler.NewUserHandler(userService)

	router := mux.NewRouter()
	routes.CreateListSubRouter(router, listHandler)
	routes.CreateUserSubrouter(router, userHandler)
	server := httptest.NewServer(router)

	return &ListTestSuite{
		DB:          db,
		Server:      server,
		ListHandler: listHandler,
		UserHandler: userHandler,
		Router:      router,
	}
}
func (s *ListTestSuite) GetServerURL() string {
	return s.Server.URL
}
func (s *ListTestSuite) GetClient() *http.Client {
	return s.Server.Client()
}

var listTestSuite *ListTestSuite
var accesstoken string

func TestAddToList(t *testing.T) {
	accessToken, _, _, err := RegisterAndLoginTestUser("list", "listadd@test.com", "1234567@M")
	require.NoError(t, err)
	accesstoken = accessToken
	body := `{
		"game_id":34126,
		"status": 2,
		"username":"list",
		"gamename":"metroid"
	}`
	client := listTestSuite.Server.Client()
	listAddRequest, err := http.NewRequest("POST", listTestSuite.Server.URL+"/list/add", strings.NewReader(body))
	require.NoError(t, err)
	assert.NotNil(t, listAddRequest)
	listAddRequest.Header.Set("Authorization", "Bearer "+accessToken)
	listAddRequest.Header.Set("Content-Type", "application/json")
	r, err := client.Do(listAddRequest)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, r.StatusCode)
}

func TestUpdateList(t *testing.T) {
	body := `{
		"game_id": 34126,
		"status": 1,
		"username":"list"
	}`
	client := listTestSuite.Server.Client()
	updateReq, err := http.NewRequest("PUT", listTestSuite.Server.URL+"/list/update", strings.NewReader(body))
	require.NoError(t, err)
	assert.NotNil(t, updateReq)
	updateReq.Header.Set("Authorization", "Bearer "+accesstoken)
	updateReq.Header.Set("Content-Type", "application/json")

	updateResp, err := client.Do(updateReq)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, updateResp.StatusCode)

}
