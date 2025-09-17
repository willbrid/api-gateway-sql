package repository

import (
	"context"

	"api-gateway-sql/internal/domain"
)

type ISQLQueryRepo interface {
	Execute(ctx context.Context, query string, params map[string]any) (*domain.SQLQueryOutput, error)
}

type Repositories struct {
	ISQLQueryRepo ISQLQueryRepo
}

func NewRepositories() *Repositories {
	return &Repositories{
		ISQLQueryRepo: NewSQLQueryRepo(),
	}
}
