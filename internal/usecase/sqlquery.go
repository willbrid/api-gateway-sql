package usecase

import (
	"api-gateway-sql/config"
	"api-gateway-sql/internal/domain"
	"api-gateway-sql/internal/repository"
)

type SQLQueryUsecase struct {
	repo   *repository.SQLQueryRepo
	config *config.Config
}

func NewSQLQueryUsecase(repo *repository.SQLQueryRepo, config *config.Config) *SQLQueryUsecase {
	return &SQLQueryUsecase{repo, config}
}

func (squ *SQLQueryUsecase) ExecuteSingle(sqlquery domain.SQLQueryInput) (*domain.SQLQueryOutput, error) {
	return nil, nil
}
