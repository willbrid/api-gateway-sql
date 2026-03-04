package usecase

import (
	"github.com/rs/zerolog"

	"github.com/willbrid/api-gateway-sql/config"
	"github.com/willbrid/api-gateway-sql/internal/dto"
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
	logger zerolog.Logger
}

func NewSQLQueryUsecase(repo *repository.SQLQueryRepo, config *config.Config, logger zerolog.Logger) *SQLQueryUsecase {
	return &SQLQueryUsecase{
		repo:   repo,
		config: config,
		logger: logger.With().Str("layer", "usecase").Str("component", "sqlquery").Logger(),
	}
}

func (squ *SQLQueryUsecase) ExecuteSingle(ctx context.Context, sqlquery *dto.SQLQueryInput) (*dto.SQLQueryOutput, error) {
	target, cfgdb, err := confighelper.GetTargetAndDatabase(squ.config, sqlquery.TargetName)
	if err != nil {
		squ.logger.Error().Err(err).Msg("unable to get target and database from config")
		return nil, err
	}

	cnx, err := external.NewDatabase(*cfgdb)
	if err != nil {
		squ.logger.Error().Err(err).Msg("unable to open database connection")
		return nil, err
	}

	squ.repo.SetDB(cnx)
	defer squ.repo.CloseDB()

	result, err := squ.repo.Execute(ctx, target.SqlQuery, sqlquery.PostParams)
	if err != nil {
		squ.logger.Error().Err(err).Msg("unable to execute single query")
		return nil, err
	}

	squ.logger.Info().Msg("single query executed")
	return result, nil
}

func (squ *SQLQueryUsecase) ExecuteInit(ctx context.Context, sqlinit *dto.SQLInitDatabaseInput) error {
	database, exist := squ.config.GetDatabaseByDataSourceName(sqlinit.Datasource)
	if !exist {
		squ.logger.Error().Msg(errUnknownDatasource.Error())
		return errUnknownDatasource
	}

	cnx, err := external.NewDatabase(database)
	if err != nil {
		squ.logger.Error().Err(err).Msg("unable to open database connection")
		return err
	}

	squ.repo.SetDB(cnx)
	defer squ.repo.CloseDB()

	queries := strings.Split(sqlinit.SQLFileContent, ";")

	if err := squ.repo.ExecuteInit(ctx, queries); err != nil {
		squ.logger.Error().Err(err).Msg("unable to execute init query")
		return err
	}

	squ.logger.Info().Msg("init query executed")
	return nil
}
