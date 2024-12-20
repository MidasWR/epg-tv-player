package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupSettingsTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}
	if err := db.AutoMigrate(&SettingsApp{}); err != nil {
		panic(err)
	}
	return db
}

func TestAddDefaultSettings(t *testing.T) {
	db := setupSettingsTestDB()
	AddDefaultSettings(db, "key1", "key2", "key3")
	var values []string
	values = append(values, GetSettings("key1", db))
	values = append(values, GetSettings("key2", db))
	values = append(values, GetSettings("key3", db))
	for _, value := range values {
		assert.Equal(t, "default", value)
	}
}

func TestAddNewSettings(t *testing.T) {
	db := setupSettingsTestDB()
	AddDefaultSettings(db, "key1", "key2")
	SetSettings("key1", "updated_value", db)
	var updatedSetting SettingsApp
	result := db.Where("key = ?", "key1").First(&updatedSetting)
	assert.NoError(t, result.Error)
	assert.Equal(t, "updated_value", updatedSetting.Value)
	var defaultSetting SettingsApp
	result = db.Where("key = ?", "key2").First(&defaultSetting)
	assert.NoError(t, result.Error)
	assert.Equal(t, "default", defaultSetting.Value)
}
