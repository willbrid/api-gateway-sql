package usecase

import (
	"api-gateway-sql/internal/domain"
	"api-gateway-sql/internal/repository"
	"api-gateway-sql/pkg/paginator"

	"context"
)

type BatchStatUsecase struct {
	repo *repository.BatchStatRepo
}

func NewBatchStatUsecase(batchStatRepo *repository.BatchStatRepo) *BatchStatUsecase {
	return &BatchStatUsecase{batchStatRepo}
}

func (b *BatchStatUsecase) ListBatchStats(ctx context.Context, pageRequest *paginator.PageRequest) (*paginator.PageResponse, error) {
	offset := pageRequest.Offset()
	limit := pageRequest.Limit()

	stats, total, err := b.repo.FindAll(ctx, offset, limit)
	if err != nil {
		return paginator.NewPageResponse(nil, 0, pageRequest), err
	}

	statsAny := make([]any, len(stats))
	for i, s := range stats {
		statsAny[i] = s
	}

	return paginator.NewPageResponse(statsAny, total, pageRequest), nil
}

func (b *BatchStatUsecase) GetBatchStatById(ctx context.Context, uid string) (*domain.BatchStat, error) {
	return b.repo.FindById(ctx, uid)
}

func (b *BatchStatUsecase) MarkCompletedBatchStat(ctx context.Context, uid string) error {
	batchStat, err := b.repo.FindById(ctx, uid)
	if err != nil {
		return err
	}

	return b.repo.UpdateLastCompleted(ctx, batchStat)
}
