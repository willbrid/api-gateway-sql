package external

import (
	"api-gateway-sql/config"

	"context"
	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type SQLServerExecutor struct {
	db *gorm.DB
}

func NewSQLServerExecutor(db config.Database) (*SQLServerExecutor, error) {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%v?database=%s&connection+timeout=%v", db.Username, db.Password, db.Host, db.Port, db.Dbname, db.Timeout.Seconds())

	cnx, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &SQLServerExecutor{db: cnx}, nil
}

func (e *SQLServerExecutor) Execute(ctx context.Context, query string, params map[string]any) (*ExecutionResult, error) {
	return executeQuery(ctx, e.db, query, params)
}
