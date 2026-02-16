package external

import (
	"github.com/willbrid/api-gateway-sql/config"

	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteDatabase struct{}

func (*SqliteDatabase) Connect(db config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s.db", db.Dbname)

	cnx, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return cnx, nil
}
