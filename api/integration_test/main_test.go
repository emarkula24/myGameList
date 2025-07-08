package integration_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	err := godotenv.Load(filepath.Join("../", ".env"))
	if err != nil {
		log.Println("No .env file found or failed to load:")
	}

	testDB, err := GetTestDataBase()
	if err != nil {
		log.Fatalf("Failed to initialize testDB: %v", err)
	}
	dbInstance := testDB.DbInstance

	userTestSuite = NewUserTestSuite(dbInstance)
	listTestSuite = NewListTestSuite(dbInstance)

	defer userTestSuite.Server.Close()
	defer listTestSuite.Server.Close()
	defer testDB.TearDown()

	os.Exit(m.Run())
}
