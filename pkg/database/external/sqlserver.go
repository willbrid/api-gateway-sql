package external

import (
	"github.com/willbrid/api-gateway-sql/config"

	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type SQLServerDatabase struct {
	db *gorm.DB
}

func (*SQLServerDatabase) Connect(db config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%v?database=%s&connection+timeout=%v", db.Username, db.Password, db.Host, db.Port, db.Dbname, db.Timeout.Seconds())

	cnx, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return cnx, nil
}
