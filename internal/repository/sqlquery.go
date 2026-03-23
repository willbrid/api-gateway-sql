package repository

import (
	"context"
	"strings"

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
	return &SQLQueryRepo{logger: logger.With().Str("layer", "repository").Str("component", "sqlqueryrepo").Logger()}
}

func (r *SQLQueryRepo) SetDB(db *gorm.DB) {
	r.db = db
}

func (r *SQLQueryRepo) CloseDB() {
	if cnx, err := r.db.DB(); err == nil {
		_ = cnx.Close()
	} else {
		r.logger.Error().Err(err).Msg("unable to close database session")
	}
}

func (r *SQLQueryRepo) ExecuteInit(ctx context.Context, sqlQueries []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, sqlQuery := range sqlQueries {
			query := strings.TrimSpace(sqlQuery)
			if query != "" {
				if err := tx.Exec(query).Error; err != nil {
					r.logger.Error().Err(err).Msg("failed to execute transaction for schema creation")
					return err
				}
			}
		}

		return nil
	})
}

func (r *SQLQueryRepo) Execute(ctx context.Context, query string, params map[string]any) (*dto.SQLQueryOutput, error) {
	parsedQuery, parsedParams := sqlqueryhelper.TransformQuery(query, params)

	if sqlqueryhelper.IsSelectQuery(parsedQuery) {
		return r.executeSelect(ctx, parsedQuery, parsedParams)
	}

	return r.executeWrite(ctx, parsedQuery, parsedParams)
}

func (r *SQLQueryRepo) executeSelect(ctx context.Context, query string, params []any) (*dto.SQLQueryOutput, error) {
	var rows []map[string]any

	if err := r.db.WithContext(ctx).Raw(query, params...).Scan(&rows).Error; err != nil {
		r.logger.Error().Err(err).Str("query", query).Msg("failed to execute select query")
		return nil, err
	}

	return &dto.SQLQueryOutput{
		Rows:         rows,
		AffectedRows: int64(len(rows)),
	}, nil
}

func (r *SQLQueryRepo) executeWrite(ctx context.Context, query string, params []any) (*dto.SQLQueryOutput, error) {
	tx := r.db.WithContext(ctx).Exec(query, params...)

	if tx.Error != nil {
		r.logger.Error().Err(tx.Error).Str("query", query).Msg("failed to execute write query")
		return nil, tx.Error
	}

	return &dto.SQLQueryOutput{
		Rows:         nil,
		AffectedRows: tx.RowsAffected,
	}, nil
}

func (r *SQLQueryRepo) ExecuteBatch(ctx context.Context, query string, params []map[string]any) error {
	cnx := r.db

	return cnx.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, param := range params {
			parsedQuery, parsedParams := sqlqueryhelper.TransformQuery(query, param)
			if err := tx.Exec(parsedQuery, parsedParams...).Error; err != nil {
				r.logger.Error().Err(tx.Error).Str("query", parsedQuery).Msg("failed to execute batch query")
				return err
			}
		}

		return nil
	})
}
