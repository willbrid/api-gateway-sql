package repository

import (
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/willbrid/api-gateway-sql/internal/dto"
	"github.com/willbrid/api-gateway-sql/internal/pkg/sqlqueryhelper"
)

type SQLQueryRepo struct {
	db     *gorm.DB
	logger zerolog.Logger
}

func NewSQLQueryRepo(logger zerolog.Logger) *SQLQueryRepo {
	return &SQLQueryRepo{logger: logger}
}

func (r *SQLQueryRepo) SetDB(db *gorm.DB) {
	r.db = db
}

func (r *SQLQueryRepo) CloseDB() {
	if cnx, err := r.db.DB(); err == nil {
		_ = cnx.Close()
	} else {
		r.logger.Error().Err(err).Str("domain", "sqlquery").Msg("unable to close database session")
	}
}

func (r *SQLQueryRepo) ExecuteInit(ctx context.Context, sqlQueries []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, sqlQuery := range sqlQueries {
			query := strings.TrimSpace(sqlQuery)
			if query != "" {
				if err := tx.Exec(query).Error; err != nil {
					r.logger.Error().Err(err).Str("domain", "sqlquery").Msg("failed to execute transaction for schema creation")
					return err
				}
			}
		}

		return nil
	})
}

func (r *SQLQueryRepo) Execute(ctx context.Context, query string, params map[string]any) (*dto.SQLQueryOutput, error) {
	var result []map[string]any
	start := time.Now()
	cnx := r.db
	parsedQuery, parsedParams := sqlqueryhelper.TransformQuery(query, params)

	if err := cnx.WithContext(ctx).Raw(parsedQuery, parsedParams...).Scan(&result).Error; err == nil {
		return &dto.SQLQueryOutput{
			Rows:         result,
			AffectedRows: int64(len(result)),
			DurationMs:   time.Since(start).Milliseconds(),
		}, nil
	}

	tx := cnx.Exec(parsedQuery, parsedParams...)
	if tx.Error != nil {
		r.logger.Error().Err(tx.Error).Str("domain", "sqlquery").Str("query", parsedQuery).Msg("failed to execute single query")
		return nil, tx.Error
	}

	return &dto.SQLQueryOutput{
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
				r.logger.Error().Err(tx.Error).Str("domain", "sqlquery").Str("query", parsedQuery).Msg("failed to execute batch query")
				return err
			}
		}

		return nil
	})
}
