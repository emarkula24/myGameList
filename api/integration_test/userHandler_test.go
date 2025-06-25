package integration_test

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"example.com/mygamelist/handler"
	"example.com/mygamelist/repository"
	"example.com/mygamelist/service"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

var dbDSN string

type TestDatabase struct {
	DbInstance *sql.DB
	container  testcontainers.Container
}

func SetupTestDatabase() *TestDatabase {

	// setup db container

	ctx := context.Background()

	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0",
		mysql.WithScripts(filepath.Join("../", "schema.sql")),
		testcontainers.WithExposedPorts("3306/tcp"),
		testcontainers.WithWaitStrategy(wait.ForExposedPort().WithStartupTimeout(300*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to create container %s", err)
	}
	// Get DSN for sql.Open
	dbDSN, err = mysqlContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to get DSN: %s", err)
	}
	db, err := sql.Open("mysql", dbDSN)
	if err != nil {
		log.Fatal("failed to setup test", err)
	}

	return &TestDatabase{
		container:  mysqlContainer,
		DbInstance: db,
	}
}
func (tdb *TestDatabase) TearDown() {
	tdb.DbInstance.Close()
	// remove test container
	if err := testcontainers.TerminateContainer(tdb.container); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
}

var testDbInstance *sql.DB

func TestMain(m *testing.M) {
	testDB := SetupTestDatabase()
	testDbInstance = testDB.DbInstance
	defer testDB.TearDown()
	os.Exit(m.Run())
}

func TestUserRegister(t *testing.T) {
	repo := repository.NewRepository(testDbInstance)
	service := service.NewUserService(repo)
	h := handler.NewHandler(service)

	body := `{
		"username": "testuser",
		"email": "testuser@example.com",
		"password": "securepassword"
	}`
	rr := httptest.NewRecorder()

	req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	assert.Equal(t, rr.Code, http.StatusOK)

	testHandler := http.HandlerFunc(h.Register)
	testHandler.ServeHTTP(rr, req)

}
