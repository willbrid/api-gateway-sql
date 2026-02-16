package external

import (
	"github.com/willbrid/api-gateway-sql/config"

	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDatabase struct{}

func (*PostgresDatabase) Connect(db config.Database) (*gorm.DB, error) {
	sslMode := "disable"
	if db.Sslmode {
		sslMode = "enable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=%v connect_timeout=%v", db.Host, db.Username, db.Password, db.Dbname, db.Port, sslMode, db.Timeout.Seconds())
	cnx, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return cnx, nil
}
