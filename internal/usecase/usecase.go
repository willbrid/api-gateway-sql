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

type ISQLBatchQueryUsecase interface {
	ExecuteBatch(ctx context.Context, sqlbatchquery *domain.SQLBatchQueryInput) error
}

type ISQLInitDatabaseUsecase interface {
	ExecuteInit(ctx context.Context, sqlinit *domain.SQLInitDatabaseInput) error
}

type Usecases struct {
	ISQLQueryUsecase        ISQLQueryUsecase
	ISQLInitDatabaseUsecase ISQLInitDatabaseUsecase
	ISQLBatchQueryUsecase   ISQLBatchQueryUsecase
}

type Deps struct {
	Repos  *repository.Repositories
	Config *config.Config
}

func NewUsecases(deps Deps) *Usecases {
	sqlQueryUsecase := NewSQLQueryUsecase(deps.Repos.ISQLQueryRepo.(*repository.SQLQueryRepo), deps.Config)
	sqlInitDatabaseUsecase := NewSQLInitDatabaseUsecase(deps.Repos.ISQLInitDatabaseRepo.(*repository.SQLInitDatabaseRepo), deps.Config)
	sqlBatchQueryUsecase := NewSQLBatchQueryUsecase(
		deps.Repos.ISQLQueryRepo.(*repository.SQLQueryRepo),
		deps.Repos.IBatchStat.(*repository.BatchStatRepo),
		deps.Repos.IBlock.(*repository.BlockRepo),
		deps.Config,
	)

	return &Usecases{
		ISQLQueryUsecase:        sqlQueryUsecase,
		ISQLInitDatabaseUsecase: sqlInitDatabaseUsecase,
		ISQLBatchQueryUsecase:   sqlBatchQueryUsecase,
	}
}
