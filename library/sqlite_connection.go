package library

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGormSqliteConnection(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
