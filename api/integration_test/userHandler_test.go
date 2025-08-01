package integration_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
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

type UserTestSuite struct {
	DB      *sql.DB
	Server  *httptest.Server
	Router  *mux.Router
	Handler *handler.UserHandler
}

func NewUserTestSuite(db *sql.DB) *UserTestSuite {
	repo := repository.NewRepository(db)
	auth := utils.AuthService{}
	svc := service.NewUserService(repo, auth)
	h := handler.NewUserHandler(svc)

	router := mux.NewRouter()
	router = routes.CreateUserSubrouter(router, h)
	server := httptest.NewServer(router)

	return &UserTestSuite{
		DB:      db,
		Server:  server,
		Handler: h,
		Router:  router,
	}
}
func (s *UserTestSuite) GetServerURL() string {
	return s.Server.URL
}

func (s *UserTestSuite) GetClient() *http.Client {
	return s.Server.Client()
}

var (
	userTestSuite *UserTestSuite
)

func TestRegister(t *testing.T) {

	body := `{
		"username": "test",
		"email":    "testregister@example.com",
		"password": "securepassword"
	}`

	// Send HTTP POST to /register
	r, err := http.Post(userTestSuite.Server.URL+"/user/register", "application/json", strings.NewReader(body))
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusCreated, r.StatusCode)

	var response struct {
		UserID int64 `json:"user_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&response)
	require.NoError(t, err)
	err = r.Body.Close()
	require.NoError(t, err)
	assert.Positive(t, response.UserID)
}

// Tests Login endpoint, depends on Register test run finishing first.
func TestLogin(t *testing.T) {
	// Same body as the register test because there is no seeded data
	body := `{
 		"username": "test",
    	"password": "securepassword"
	}`

	r, err := http.Post(userTestSuite.Server.URL+"/user/login", "application/json", strings.NewReader(body))
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	var response struct {
		AccessToken string `json:"accessToken"`
		UserID      int    `json:"userId"`
	}
	err = json.NewDecoder(r.Body).Decode(&response)
	require.NoError(t, err)
	err = r.Body.Close()
	require.NoError(t, err)
	assert.NotEmpty(t, response.AccessToken)

}

func TestRefresh(t *testing.T) {
	registerBody := `{
		"username": "testrefresh",
		"email":    "testrefresh@example.com",
		"password": "123456778@M"
	}`

	loginBody := `{
 		"username": "testrefresh",
		"password": "123456778@M"
	}`

	r, err := http.Post(userTestSuite.Server.URL+"/user/register", "application/json", strings.NewReader(registerBody))
	require.NoError(t, err)
	require.NotNil(t, r)
	err = r.Body.Close()
	require.NoError(t, err)
	r, err = http.Post(userTestSuite.Server.URL+"/user/login", "application/json", strings.NewReader(loginBody))
	require.NoError(t, err)
	require.NotNil(t, r)

	var response struct {
		AccessToken string `json:"accessToken"`
		UserId      int    `json:"userId"`
	}

	cookies := r.Cookies()
	err = json.NewDecoder(r.Body).Decode(&response)
	require.NoError(t, err)
	err = r.Body.Close()
	require.NoError(t, err)

	refreshData := map[string]any{
		"username": "testrefresh",
		"userId":   response.UserId,
	}

	refreshBodyBytes, err := json.Marshal(refreshData)
	require.NoError(t, err)
	client := userTestSuite.Server.Client()
	refreshReq, err := http.NewRequest("POST", userTestSuite.Server.URL+"/user/refresh", bytes.NewReader(refreshBodyBytes))
	require.NoError(t, err)
	refreshReq.Header.Set("Content-Type", "application/json")

	for _, cookie := range cookies {
		refreshReq.AddCookie(cookie)
	}

	r, err = client.Do(refreshReq)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	log.Println(r.Header)

	var refreshResponse struct {
		AccessTokenFromRefresh string `json:"accessToken"`
	}
	err = json.NewDecoder(r.Body).Decode(&refreshResponse)
	require.NoError(t, err)
	err = r.Body.Close()
	require.NoError(t, err)
	assert.NotEmpty(t, refreshResponse.AccessTokenFromRefresh)

}
