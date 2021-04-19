package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewTestDatabase(t *testing.T) *gorm.DB {
	var cfg struct {
		Database struct {
			User     string
			Password string
			Name     string
			Port     string
			Host     string
			SSLMode  string
			Timezone string
		}
	}
	cfg.Database.Host = "0.0.0.0"
	cfg.Database.User = "gocallbacksvc"
	cfg.Database.Password = "secret1234"
	cfg.Database.Name = "gocallbacksvc_test"
	cfg.Database.Port = "5433"
	cfg.Database.SSLMode = "disable"
	cfg.Database.Timezone = "Europe/Madrid"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode,
		cfg.Database.Timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(t, err, "opening database connection through dsl")

	db.AutoMigrate(Callback{})

	return db
}

func CleanupTestDatabase(gdb *gorm.DB) {
	gdb.Exec("DROP SCHEMA public CASCADE")
	gdb.Exec("CREATE SCHEMA public")
	gdb.Migrator().CreateTable(&Callback{})
}
