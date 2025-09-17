package repository

import (
	"context"

	"api-gateway-sql/internal/domain"
)

type ISQLQuery interface {
	Execute(ctx context.Context, query string, params map[string]any) (*domain.SQLQueryOutput, error)
}
