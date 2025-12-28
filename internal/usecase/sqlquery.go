package usecase

import (
	"github.com/willbrid/api-gateway-sql/config"
	"github.com/willbrid/api-gateway-sql/internal/domain"
	"github.com/willbrid/api-gateway-sql/internal/pkg/confighelper"
	"github.com/willbrid/api-gateway-sql/internal/repository"
	"github.com/willbrid/api-gateway-sql/pkg/database/external"

	"context"
	"errors"
	"strings"
)

var (
	errUnknownDatasource error = errors.New("unknown datasource name")
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

func (squ *SQLQueryUsecase) ExecuteInit(ctx context.Context, sqlinit *domain.SQLInitDatabaseInput) error {
	database, exist := squ.config.GetDatabaseByDataSourceName(sqlinit.Datasource)
	if !exist {
		return errUnknownDatasource
	}

	cnx, err := external.NewDatabase(database)
	if err != nil {
		return err
	}

	squ.repo.SetDB(cnx)
	defer squ.repo.CloseDB()

	queries := strings.Split(sqlinit.SQLFileContent, ";")

	return squ.repo.ExecuteInit(ctx, queries)
}
