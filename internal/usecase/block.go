package usecase

import (
	"github.com/rs/zerolog"

	"github.com/willbrid/api-gateway-sql/internal/domain"
	"github.com/willbrid/api-gateway-sql/internal/dto/paginator"
	"github.com/willbrid/api-gateway-sql/internal/repository"

	"context"
)

type BlockUsecase struct {
	repo   *repository.BlockRepo
	logger zerolog.Logger
}

func NewBlockUsecase(blockRepo *repository.BlockRepo, logger zerolog.Logger) *BlockUsecase {
	return &BlockUsecase{
		repo:   blockRepo,
		logger: logger.With().Str("layer", "usecase").Str("component", "block").Logger(),
	}
}

func (bu *BlockUsecase) ListBlocksByBatchStat(ctx context.Context, batchStatId string, pageRequest *paginator.PageRequest) (*paginator.PageResponse, error) {
	offset := pageRequest.Offset()
	limit := pageRequest.Limit()

	blocks, total, err := bu.repo.FindAllByBatchStatID(ctx, batchStatId, offset, limit)
	if err != nil {
		bu.logger.Error().Err(err).Msg("failed to get blocks")
		return paginator.NewPageResponse(nil, 0, pageRequest), err
	}

	blocksAny := make([]any, len(blocks))
	for i, b := range blocks {
		blocksAny[i] = b
	}

	bu.logger.Info().Msg("blocks list got")
	return paginator.NewPageResponse(blocksAny, total, pageRequest), nil
}

func (bu *BlockUsecase) GetBlockById(ctx context.Context, uid string) (*domain.Block, error) {
	block, err := bu.repo.FindById(ctx, uid)

	if err != nil {
		bu.logger.Error().Err(err).Str("block_id", uid).Msg("failed to get block by id")
		return nil, err
	}

	bu.logger.Info().Str("block_id", uid).Msg("block found")
	return block, nil
}
