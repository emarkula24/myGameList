package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

var dbDSN string

func TestMain(m *testing.M) {
	ctx := context.Background()

	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0",
		mysql.WithScripts(filepath.Join("../", "schema.sql")),
		testcontainers.WithExposedPorts("3306/tcp"),
		testcontainers.WithWaitStrategy(wait.ForExposedPort().WithStartupTimeout(120*time.Second)),
	)

	defer func() {
		if err := testcontainers.TerminateContainer(mysqlContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}
	// Get DSN for sql.Open
	dbDSN, err = mysqlContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to get DSN: %s", err)
	}
	// Run all tests in the package
	exitCode := m.Run()

	// Exit with the appropriate status code
	os.Exit(exitCode)
}

func TestUserRegister(t *testing.T) {
	db, err := sql.Open("mysql", dbDSN)
	require.NoError(t, err)

	var name string
	err = db.QueryRow(`SELECT name FROM profile`).Scan(&name)
	assert.Nil(t, err)
	defer db.Close()

	fmt.Println("success")

}
