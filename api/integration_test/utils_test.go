package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabase struct {
	DbInstance *sql.DB
	container  testcontainers.Container
}

var (
	once    sync.Once
	testDB  *TestDatabase
	initErr error
	dbDSN   string
)

func SetupTestDatabase() {

	// setup db container

	ctx := context.Background()

	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0",
		mysql.WithScripts(filepath.Join("../", "schema.sql")),
		testcontainers.WithExposedPorts("3360/tcp"),
		testcontainers.WithWaitStrategy(wait.ForExposedPort().WithStartupTimeout(160*time.Second)),
	)
	if err != nil {
		initErr = err
		log.Fatalf("failed to create container %s", err)
		return
	}
	// Get DSN for sql.Open
	dbDSN, err = mysqlContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to get DSN: %s", err)
		initErr = err
		return
	}

	db, err := sql.Open("mysql", dbDSN)
	if err != nil {
		log.Fatal("failed to setup test", err)
		initErr = err
		return
	}
	db.SetConnMaxIdleTime(time.Minute * 5)
	db.SetMaxIdleConns(0)
	testDB = &TestDatabase{
		container:  mysqlContainer,
		DbInstance: db,
	}
}
func GetTestDataBase() (*TestDatabase, error) {
	once.Do(SetupTestDatabase)
	return testDB, initErr
}

func (tdb *TestDatabase) TearDown() {
	err := tdb.DbInstance.Close()
	if err != nil {
		log.Printf("failed to close dbinrance: %s", err)
	}
	// remove test container
	if err := testcontainers.TerminateContainer(tdb.container); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
}

func RegisterAndLoginTestUser(
	username, email, password string) (accessToken string, userId string, Username string, err error) {

	registerData := map[string]any{
		"username": username,
		"email":    email,
		"password": password,
	}
	registerBody, err := json.Marshal(registerData)
	if err != nil {
		return "", "", "", fmt.Errorf("marshal register data: %w", err)
	}

	registerResp, err := userTestSuite.GetClient().Post(userTestSuite.GetServerURL()+"/user/register", "application/json", bytes.NewReader(registerBody))
	if err != nil {
		return "", "", "", fmt.Errorf("register failed: %w", err)
	}

	if registerResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(registerResp.Body)
		return "", "", "", fmt.Errorf("register bad status: %d, body: %s", registerResp.StatusCode, string(body))
	}

	err = registerResp.Body.Close()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to close body: %w", err)
	}

	loginData := map[string]any{
		"username": username,
		"password": password,
	}
	loginBody, err := json.Marshal(loginData)
	if err != nil {
		return "", "", "", fmt.Errorf("marshal login data: %w", err)
	}

	loginResp, err := userTestSuite.GetClient().Post(userTestSuite.GetServerURL()+"/user/login", "application/json", bytes.NewReader(loginBody))
	if err != nil {
		return "", "", "", fmt.Errorf("login failed: %w", err)
	}

	if loginResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(loginResp.Body)
		return "", "", "", fmt.Errorf("login bad status: %d, body: %s", loginResp.StatusCode, string(body))
	}

	var response struct {
		AccessToken string `json:"accessToken"`
		UserId      string `json:"userId"`
		UserName    string `json:"username"`
	}
	if err := json.NewDecoder(loginResp.Body).Decode(&response); err != nil {
		return "", "", "", fmt.Errorf("decode login response: %w", err)
	}

	err = loginResp.Body.Close()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to close body: %w", err)
	}

	return response.AccessToken, response.UserId, response.UserName, nil
}
