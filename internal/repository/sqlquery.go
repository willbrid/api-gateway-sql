package repository

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"

	"api-gateway-sql/internal/domain"
	"api-gateway-sql/internal/pkg/sqlqueryhelper"
)

type SQLQueryRepo struct {
	db *gorm.DB
}

func NewSQLQueryRepo() *SQLQueryRepo {
	return &SQLQueryRepo{}
}

func (r *SQLQueryRepo) SetDB(db *gorm.DB) {
	r.db = db
}

func (r *SQLQueryRepo) CloseDB() {
	if cnx, err := r.db.DB(); err == nil {
		cnx.Close()
	}
}

func (r *SQLQueryRepo) ExecuteInit(ctx context.Context, sqlQueries []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, sqlQuery := range sqlQueries {
			query := strings.TrimSpace(sqlQuery)
			if query != "" {
				if err := tx.Exec(query).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *SQLQueryRepo) Execute(ctx context.Context, query string, params map[string]any) (*domain.SQLQueryOutput, error) {
	var result []map[string]any
	start := time.Now()
	cnx := r.db
	parsedQuery, parsedParams := sqlqueryhelper.TransformQuery(query, params)

	if err := cnx.WithContext(ctx).Raw(parsedQuery, parsedParams...).Scan(&result).Error; err == nil {
		return &domain.SQLQueryOutput{
			Rows:         result,
			AffectedRows: int64(len(result)),
			DurationMs:   time.Since(start).Milliseconds(),
		}, nil
	}

	tx := cnx.Exec(parsedQuery, parsedParams...)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &domain.SQLQueryOutput{
		Rows:         nil,
		AffectedRows: tx.RowsAffected,
		DurationMs:   time.Since(start).Milliseconds(),
	}, nil
}

func (r *SQLQueryRepo) ExecuteBatch(ctx context.Context, query string, params []map[string]any) error {
	cnx := r.db

	return cnx.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, param := range params {
			parsedQuery, parsedParams := sqlqueryhelper.TransformQuery(query, param)
			if err := tx.Exec(parsedQuery, parsedParams...).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
