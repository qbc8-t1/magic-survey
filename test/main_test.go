package api

import (
	"log"
	"os"
	"testing"

	"github.com/QBC8-Team1/magic-survey/pkg/db"
	"gorm.io/gorm"
)

var (
	testDB *gorm.DB
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = db.CreateTestDatabase()
	if err != nil {
		log.Fatalf("Failed to create test database: %v", err)
	}

	code := m.Run()

	_ = db.CloseTestDatabase(testDB)
	err = os.Remove("./test_db.db") // Replace with the actual path to test.db
	if err != nil {
		log.Fatalf("Failed to delete test database file: %v", err)
	}

	os.Exit(code)
}
