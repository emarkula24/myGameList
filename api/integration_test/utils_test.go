package integration_test

import (
	"context"
	"database/sql"
	"log"
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
		testcontainers.WithExposedPorts("3306/tcp"),
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
	tdb.DbInstance.Close()
	// remove test container
	if err := testcontainers.TerminateContainer(tdb.container); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
}
