package integration

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

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
	a := os.Getenv("TESTCONTAINER_ADDRESS")
	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0",
		mysql.WithScripts(filepath.Join("../", "schema.sql")),
		testcontainers.WithExposedPorts(a),
		testcontainers.WithWaitStrategy(wait.ForExposedPort().WithStartupTimeout(120*time.Second)),
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
