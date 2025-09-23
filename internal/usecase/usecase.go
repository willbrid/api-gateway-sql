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

type ISQLInitDatabaseUsecase interface {
	ExecuteInit(ctx context.Context, sqlinit *domain.SQLInitDatabaseInput) error
}

type Usecases struct {
	ISQLQueryUsecase        ISQLQueryUsecase
	ISQLInitDatabaseUsecase ISQLInitDatabaseUsecase
}

type Deps struct {
	Repos  *repository.Repositories
	Config *config.Config
}

func NewUsecases(deps Deps) *Usecases {
	sqlQueryUsecase := NewSQLQueryUsecase(deps.Repos.ISQLQueryRepo.(*repository.SQLQueryRepo), deps.Config)
	sqlInitDatabaseUsecase := NewSQLInitDatabaseUsecase(deps.Repos.ISQLInitDatabaseRepo.(*repository.SQLInitDatabaseRepo), deps.Config)

	return &Usecases{
		ISQLQueryUsecase:        sqlQueryUsecase,
		ISQLInitDatabaseUsecase: sqlInitDatabaseUsecase,
	}
}
