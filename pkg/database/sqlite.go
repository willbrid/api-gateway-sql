package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteAppDatabase struct {
	Db *gorm.DB
}

func NewSqliteAppDatabase(sqlitedb string) *SqliteAppDatabase {
	dsn := fmt.Sprintf("/data/%s.db", sqlitedb)

	cnx, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil
	}

	return &SqliteAppDatabase{Db: cnx}
}
