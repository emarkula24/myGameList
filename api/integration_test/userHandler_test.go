package integration

import (
	"database/sql"
	"encoding/json"
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
		"username": "testuser",
		"email":    "testuser@example.com",
		"password": "securepassword"
	}`

	// Send HTTP POST to /register
	r, err := http.Post(userTestSuite.Server.URL+"/register", "application/json", strings.NewReader(body))
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
 		"username": "testuser",
    	"password": "securepassword"
	}`

	r, err := http.Post(userTestSuite.Server.URL+"/login", "application/json", strings.NewReader(body))
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
	body := `{
 		"username": "testuser"
	}`

}
