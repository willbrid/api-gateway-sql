package usecase

import (
	"github.com/willbrid/api-gateway-sql/internal/domain"
	"github.com/willbrid/api-gateway-sql/internal/repository"
	"github.com/willbrid/api-gateway-sql/pkg/paginator"

	"context"
)

type BlockUsecase struct {
	repo *repository.BlockRepo
}

func NewBlockUsecase(blockRepo *repository.BlockRepo) *BlockUsecase {
	return &BlockUsecase{blockRepo}
}

func (bu *BlockUsecase) ListBlocksByBatchStat(ctx context.Context, batchStatId string, pageRequest *paginator.PageRequest) (*paginator.PageResponse, error) {
	offset := pageRequest.Offset()
	limit := pageRequest.Limit()

	blocks, total, err := bu.repo.FindAllByBatchStatID(ctx, batchStatId, offset, limit)
	if err != nil {
		return paginator.NewPageResponse(nil, 0, pageRequest), err
	}

	blocksAny := make([]any, len(blocks))
	for i, b := range blocks {
		blocksAny[i] = b
	}

	return paginator.NewPageResponse(blocksAny, total, pageRequest), nil
}

func (bu *BlockUsecase) GetBlockById(ctx context.Context, uid string) (*domain.Block, error) {
	return bu.repo.FindById(ctx, uid)
}
