package repository

import (
	"context"

	"api-gateway-sql/internal/domain"

	"gorm.io/gorm"
)

type ISQLQueryRepo interface {
	Execute(ctx context.Context, query string, params map[string]any) (*domain.SQLQueryOutput, error)
	ExecuteBatch(ctx context.Context, query string, params []map[string]any) error
	ExecuteInit(ctx context.Context, sqlQueries []string) error
}

type IBatchStat interface {
	Create(ctx context.Context, targetName string) (*domain.BatchStat, error)
	UpdateLastCompleted(ctx context.Context, batchStat *domain.BatchStat) error
	AddBlockToBatchStat(ctx context.Context, bs *domain.BatchStat, block *domain.Block) (*domain.Block, error)
	FindAll(ctx context.Context, offset, limit int) ([]*domain.BatchStat, int64, error)
	FindById(ctx context.Context, uid string) (*domain.BatchStat, error)
	CountUncompletedBatchStat(ctx context.Context) (int64, error)
}

type IBlock interface {
	Update(ctx context.Context, block *domain.Block, failureRange *domain.FailureRange, isSuccess bool) error
	FindAllByBatchStatID(ctx context.Context, batchStatId string, offset, limit int) ([]*domain.Block, int64, error)
	FindById(ctx context.Context, uid string) (*domain.Block, error)
}

type Repositories struct {
	ISQLQueryRepo ISQLQueryRepo
	IBatchStat    IBatchStat
	IBlock        IBlock
}

func NewRepositories(appDb *gorm.DB) *Repositories {
	return &Repositories{
		ISQLQueryRepo: NewSQLQueryRepo(),
		IBatchStat:    NewBatchStatRepo(appDb),
		IBlock:        NewBlockRepo(appDb),
	}
}
