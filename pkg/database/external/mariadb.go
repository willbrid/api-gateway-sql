package external

import (
	"api-gateway-sql/config"

	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MariadbExecutor struct {
	db *gorm.DB
}

func NewMariadbExecutor(db config.Database) (*MariadbExecutor, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%vs", db.Username, db.Password, db.Host, db.Port, db.Dbname, db.Timeout.Seconds())

	cnx, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &MariadbExecutor{db: cnx}, nil
}

func (e *MariadbExecutor) Execute(ctx context.Context, query string, params map[string]any) (*ExecutionResult, error) {
	return executeQuery(ctx, e.db, query, params)
}
