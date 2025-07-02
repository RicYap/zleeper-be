// tests/order_item_test.go
package tests

import (
	"testing"

	"zleeper-be/config"
	"zleeper-be/internal/models"
	"zleeper-be/pkg/database"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := database.InitDB(config.DBConfig{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "password",
		DBName:   "test_db",
	})
	assert.NoError(t, err)

	// Clean up and migrate
	db.Migrator().DropTable(&models.OrderItem{})
	db.AutoMigrate(&models.OrderItem{})

	return db
}