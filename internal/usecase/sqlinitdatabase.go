package usecase

import (
	"api-gateway-sql/config"
	"api-gateway-sql/internal/domain"
	"api-gateway-sql/internal/repository"
	"api-gateway-sql/pkg/database/external"

	"context"
	"errors"
	"strings"
)

var (
	errUnknownDatasource error = errors.New("unknown datasource name")
)

type SQLInitDatabaseUsecase struct {
	repo   *repository.SQLInitDatabaseRepo
	config *config.Config
}

func NewSQLInitDatabaseUsecase(repo *repository.SQLInitDatabaseRepo, config *config.Config) *SQLInitDatabaseUsecase {
	return &SQLInitDatabaseUsecase{repo, config}
}

func (sid *SQLInitDatabaseUsecase) ExecuteInit(ctx context.Context, sqlinit *domain.SQLInitDatabaseInput) error {
	database, exist := sid.config.GetDatabaseByDataSourceName(sqlinit.Datasource)
	if !exist {
		return errUnknownDatasource
	}

	cnx, err := external.NewDatabase(database)
	if err != nil {
		return err
	}

	sid.repo.SetDB(cnx)
	defer sid.repo.CloseDB()

	queries := strings.Split(sqlinit.SQLFileContent, ";")

	return sid.repo.ExecuteInit(ctx, queries)
}
