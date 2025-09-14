package external

import (
	"api-gateway-sql/config"

	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLExecutor struct {
	db *gorm.DB
}

func NewMySQLExecutor(db config.Database) (*MySQLExecutor, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%vs", db.Username, db.Password, db.Host, db.Port, db.Dbname, db.Timeout.Seconds())

	cnx, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &MySQLExecutor{db: cnx}, nil
}

func (e *MySQLExecutor) Execute(ctx context.Context, query string, params map[string]any) (*ExecutionResult, error) {
	return executeQuery(ctx, e.db, query, params)
}
