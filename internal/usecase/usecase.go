package usecase

import (
	"github.com/rs/zerolog"

	"github.com/willbrid/api-gateway-sql/config"
	"github.com/willbrid/api-gateway-sql/internal/domain"
	"github.com/willbrid/api-gateway-sql/internal/dto"
	"github.com/willbrid/api-gateway-sql/internal/dto/paginator"
	"github.com/willbrid/api-gateway-sql/internal/repository"

	"context"
)

type ISQLQueryUsecase interface {
	ExecuteSingle(ctx context.Context, sqlquery *dto.SQLQueryInput) (*dto.SQLQueryOutput, error)
	ExecuteInit(ctx context.Context, sqlinit *dto.SQLInitDatabaseInput) error
}

type ISQLBatchQueryUsecase interface {
	ExecuteBatch(ctx context.Context, sqlbatchquery *dto.SQLBatchQueryInput) error
}

type IBatchStatUsecase interface {
	ListBatchStats(ctx context.Context, pageRequest *paginator.PageRequest) (*paginator.PageResponse, error)
	GetBatchStatById(ctx context.Context, uid string) (*domain.BatchStat, error)
	MarkCompletedBatchStat(ctx context.Context, uid string) error
	IsAllBatchStatClosed(ctx context.Context) (bool, error)
}

type IBlockUsecase interface {
	ListBlocksByBatchStat(ctx context.Context, batchStatId string, pageRequest *paginator.PageRequest) (*paginator.PageResponse, error)
	GetBlockById(ctx context.Context, uid string) (*domain.Block, error)
}

type Usecases struct {
	ISQLQueryUsecase      ISQLQueryUsecase
	ISQLBatchQueryUsecase ISQLBatchQueryUsecase
	IBatchStatUsecase     IBatchStatUsecase
	IBlockUsecase         IBlockUsecase
}

type Deps struct {
	Repos  *repository.Repositories
	Config *config.Config
	Logger zerolog.Logger
}

func NewUsecases(deps Deps) *Usecases {
	sqlQueryUsecase := NewSQLQueryUsecase(deps.Repos.ISQLQueryRepo.(*repository.SQLQueryRepo), deps.Config, deps.Logger)
	sqlBatchQueryUsecase := NewSQLBatchQueryUsecase(
		deps.Repos.ISQLQueryRepo.(*repository.SQLQueryRepo),
		deps.Repos.IBatchStat.(*repository.BatchStatRepo),
		deps.Repos.IBlock.(*repository.BlockRepo),
		deps.Config,
		deps.Logger,
	)
	batchStatUsecase := NewBatchStatUsecase(deps.Repos.IBatchStat.(*repository.BatchStatRepo), deps.Logger)
	blockUsecase := NewBlockUsecase(deps.Repos.IBlock.(*repository.BlockRepo), deps.Logger)

	return &Usecases{
		ISQLQueryUsecase:      sqlQueryUsecase,
		ISQLBatchQueryUsecase: sqlBatchQueryUsecase,
		IBatchStatUsecase:     batchStatUsecase,
		IBlockUsecase:         blockUsecase,
	}
}
