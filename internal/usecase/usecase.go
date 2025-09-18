package usecase

import (
	"api-gateway-sql/config"
	"api-gateway-sql/internal/domain"
	"api-gateway-sql/internal/repository"

	"context"
)

type ISQLQueryUsecase interface {
	ExecuteSingle(ctx context.Context, sqlquery *domain.SQLQueryInput) (*domain.SQLQueryOutput, error)
}

type Usecases struct {
	ISQLQueryUsecase ISQLQueryUsecase
}

type Deps struct {
	Repos *repository.Repositories
}

func NewUsecases(deps Deps, config *config.Config) *Usecases {
	sqlQueryUsecase := NewSQLQueryUsecase(deps.Repos.ISQLQueryRepo.(*repository.SQLQueryRepo), config)

	return &Usecases{
		ISQLQueryUsecase: sqlQueryUsecase,
	}
}
