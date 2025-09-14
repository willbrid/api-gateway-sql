package external

import (
	"api-gateway-sql/config"

	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresExecutor struct {
	db *gorm.DB
}

func NewPostgresExecutor(db config.Database) (*PostgresExecutor, error) {
	sslMode := "disable"
	if db.Sslmode {
		sslMode = "enable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=%v connect_timeout=%v", db.Host, db.Username, db.Password, db.Dbname, db.Port, sslMode, db.Timeout.Seconds())
	cnx, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &PostgresExecutor{db: cnx}, nil
}

func (e *PostgresExecutor) Execute(ctx context.Context, query string, params map[string]any) (*ExecutionResult, error) {
	return executeQuery(ctx, e.db, query, params)
}
