package usecase

import (
	"github.com/rs/zerolog"

	"github.com/willbrid/api-gateway-sql/internal/domain"
	"github.com/willbrid/api-gateway-sql/internal/dto/paginator"
	"github.com/willbrid/api-gateway-sql/internal/repository"

	"context"
)

type BatchStatUsecase struct {
	repo   *repository.BatchStatRepo
	logger zerolog.Logger
}

func NewBatchStatUsecase(batchStatRepo *repository.BatchStatRepo, logger zerolog.Logger) *BatchStatUsecase {
	return &BatchStatUsecase{
		repo:   batchStatRepo,
		logger: logger.With().Str("layer", "usecase").Str("component", "batchstat").Logger(),
	}
}

func (b *BatchStatUsecase) IsAllBatchStatClosed(ctx context.Context) (bool, error) {
	totalUncompletedBatchStat, err := b.repo.CountUncompletedBatchStat(ctx)
	if err != nil {
		b.logger.Error().Err(err).Msg("failed to get uncompleted batchstat")
		return false, err
	}

	b.logger.Info().Msg("uncompleted batchstat count retrieved")
	return totalUncompletedBatchStat <= 0, nil
}

func (b *BatchStatUsecase) ListBatchStats(ctx context.Context, pageRequest *paginator.PageRequest) (*paginator.PageResponse, error) {
	offset := pageRequest.Offset()
	limit := pageRequest.Limit()

	stats, total, err := b.repo.FindAll(ctx, offset, limit)
	if err != nil {
		b.logger.Error().Err(err).Msg("failed to get batchstat list")
		return paginator.NewPageResponse(nil, 0, pageRequest), err
	}

	statsAny := make([]any, len(stats))
	for i, s := range stats {
		statsAny[i] = s
	}

	b.logger.Info().Msg("batchstat list got")
	return paginator.NewPageResponse(statsAny, total, pageRequest), nil
}

func (b *BatchStatUsecase) GetBatchStatById(ctx context.Context, uid string) (*domain.BatchStat, error) {
	batchStat, err := b.repo.FindById(ctx, uid)

	if err != nil {
		b.logger.Error().Err(err).Str("batchstat_id", uid).Msg("failed to get batchstat by id")
		return nil, err
	}

	b.logger.Info().Str("batchstat_id", uid).Msg("batchstat found")
	return batchStat, nil
}

func (b *BatchStatUsecase) MarkCompletedBatchStat(ctx context.Context, uid string) error {
	batchStat, err := b.repo.FindById(ctx, uid)
	if err != nil {
		b.logger.Error().Err(err).Str("batchstat_id", uid).Msg("failed to mark batchstat completed")
		return err
	}

	b.logger.Info().Msg("batchstat marked completed")
	return b.repo.UpdateLastCompleted(ctx, batchStat)
}
