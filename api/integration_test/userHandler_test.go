package integration

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"example.com/mygamelist/handler"
	"example.com/mygamelist/repository"
	"example.com/mygamelist/route"
	"example.com/mygamelist/service"
	"example.com/mygamelist/utils"
	_ "github.com/go-sql-driver/mysql"
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
	router = route.CreateUserSubrouter(router, h)
	server := httptest.NewServer(router)

	return &UserTestSuite{
		DB:      db,
		Server:  server,
		Handler: h,
		Router:  router,
	}
}

func (ts *UserTestSuite) Close() {
	ts.Server.Close()
}

var (
	testDbInstance *sql.DB
	userTestSuite  *UserTestSuite
)

func TestMain(m *testing.M) {
	testDB := SetupTestDatabase()
	testDbInstance = testDB.DbInstance
	userTestSuite = NewUserTestSuite(testDbInstance)

	defer userTestSuite.Close()
	defer testDB.TearDown()

	os.Exit(m.Run())
}

func TestRegister(t *testing.T) {

	body := `{
		"username": "test",
		"email":    "testregister@example.com",
		"password": "securepassword"
	}`

	// Send HTTP POST to /register
	r, err := http.Post(userTestSuite.Server.URL+"/user/register", "application/json", strings.NewReader(body))
	require.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	defer r.Body.Close()
	var response struct {
		UserID int64 `json:"user_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&response)
	require.Nil(t, err)
	assert.Greater(t, response.UserID, int64(0))
}

// Tests Login endpoint, depends on Register test run finishing first.
func TestLogin(t *testing.T) {
	// Same body as the register test because there is no seeded data
	body := `{
 		"username": "test",
    	"password": "securepassword"
	}`

	r, err := http.Post(userTestSuite.Server.URL+"/user/login", "application/json", strings.NewReader(body))
	require.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	defer r.Body.Close()

	var response struct {
		AccessToken string `json:"accessToken"`
	}
	err = json.NewDecoder(r.Body).Decode(&response)
	require.Nil(t, err)
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

	refreshBody := `{
		"username": "testrefresh"
	}`

	r, err := http.Post(userTestSuite.Server.URL+"/user/register", "application/json", strings.NewReader(registerBody))
	require.Nil(t, err)
	require.NotNil(t, r)
	r.Body.Close()
	r, err = http.Post(userTestSuite.Server.URL+"/user/login", "application/json", strings.NewReader(loginBody))
	require.Nil(t, err)
	require.NotNil(t, r)

	var response struct {
		AccessToken string `json:"accessToken"`
	}

	cookies := r.Cookies()
	err = json.NewDecoder(r.Body).Decode(&response)
	require.Nil(t, err)
	r.Body.Close()

	client := userTestSuite.Server.Client()
	refreshReq, err := http.NewRequest("POST", userTestSuite.Server.URL+"/user/refresh", strings.NewReader(refreshBody))
	require.Nil(t, err)
	refreshReq.Header.Set("Content-Type", "application/json")

	for _, cookie := range cookies {
		refreshReq.AddCookie(cookie)
	}

	r, err = client.Do(refreshReq)
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	log.Println(r.Header)
	defer r.Body.Close()

	var refreshResponse struct {
		AccessTokenFromRefresh string `json:"accessToken"`
	}
	err = json.NewDecoder(r.Body).Decode(&refreshResponse)
	require.Nil(t, err)
	assert.NotEmpty(t, refreshResponse.AccessTokenFromRefresh)

}
