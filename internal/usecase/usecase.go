package usecase

import (
	"github.com/willbrid/api-gateway-sql/config"
	"github.com/willbrid/api-gateway-sql/internal/domain"
	"github.com/willbrid/api-gateway-sql/internal/repository"
	"github.com/willbrid/api-gateway-sql/pkg/paginator"

	"context"
)

type ISQLQueryUsecase interface {
	ExecuteSingle(ctx context.Context, sqlquery *domain.SQLQueryInput) (*domain.SQLQueryOutput, error)
	ExecuteInit(ctx context.Context, sqlinit *domain.SQLInitDatabaseInput) error
}

type ISQLBatchQueryUsecase interface {
	ExecuteBatch(ctx context.Context, sqlbatchquery *domain.SQLBatchQueryInput) error
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
}

func NewUsecases(deps Deps) *Usecases {
	sqlQueryUsecase := NewSQLQueryUsecase(deps.Repos.ISQLQueryRepo.(*repository.SQLQueryRepo), deps.Config)
	sqlBatchQueryUsecase := NewSQLBatchQueryUsecase(
		deps.Repos.ISQLQueryRepo.(*repository.SQLQueryRepo),
		deps.Repos.IBatchStat.(*repository.BatchStatRepo),
		deps.Repos.IBlock.(*repository.BlockRepo),
		deps.Config,
	)
	batchStatUsecase := NewBatchStatUsecase(deps.Repos.IBatchStat.(*repository.BatchStatRepo))
	blockUsecase := NewBlockUsecase(deps.Repos.IBlock.(*repository.BlockRepo))

	return &Usecases{
		ISQLQueryUsecase:      sqlQueryUsecase,
		ISQLBatchQueryUsecase: sqlBatchQueryUsecase,
		IBatchStatUsecase:     batchStatUsecase,
		IBlockUsecase:         blockUsecase,
	}
}
