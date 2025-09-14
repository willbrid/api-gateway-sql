package external

import (
	"api-gateway-sql/config"

	"context"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteExecutor struct {
	db *gorm.DB
}

func NewSqliteExecutor(db config.Database) (*SqliteExecutor, error) {
	dsn := fmt.Sprintf("%s.db", db.Dbname)

	cnx, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &SqliteExecutor{db: cnx}, nil
}

func (e *SqliteExecutor) Execute(ctx context.Context, query string, params map[string]any) (*ExecutionResult, error) {
	return executeQuery(ctx, e.db, query, params)
}
