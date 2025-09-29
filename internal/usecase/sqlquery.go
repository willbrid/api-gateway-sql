package usecase

import (
	"api-gateway-sql/config"
	"api-gateway-sql/internal/domain"
	"api-gateway-sql/internal/pkg/confighelper"
	"api-gateway-sql/internal/repository"
	"api-gateway-sql/pkg/database/external"

	"context"
)

type SQLQueryUsecase struct {
	repo   *repository.SQLQueryRepo
	config *config.Config
}

func NewSQLQueryUsecase(repo *repository.SQLQueryRepo, config *config.Config) *SQLQueryUsecase {
	return &SQLQueryUsecase{repo, config}
}

func (squ *SQLQueryUsecase) ExecuteSingle(ctx context.Context, sqlquery *domain.SQLQueryInput) (*domain.SQLQueryOutput, error) {
	target, cfgdb, err := confighelper.GetTargetAndDatabase(squ.config, sqlquery.TargetName)
	if err != nil {
		return nil, err
	}

	cnx, err := external.NewDatabase(*cfgdb)
	if err != nil {
		return nil, err
	}

	squ.repo.SetDB(cnx)
	defer squ.repo.CloseDB()

	result, err := squ.repo.Execute(ctx, target.SqlQuery, sqlquery.PostParams)
	if err != nil {
		return nil, err
	}

	return result, nil
}
