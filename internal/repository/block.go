package repository

import (
	"github.com/willbrid/api-gateway-sql/internal/domain"

	"context"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type BlockRepo struct {
	appDb  *gorm.DB
	logger zerolog.Logger
}

func NewBlockRepo(appDb *gorm.DB, logger zerolog.Logger) *BlockRepo {
	return &BlockRepo{appDb, logger}
}

func (b *BlockRepo) Update(ctx context.Context, block *domain.Block, failureRange *domain.FailureRange, isSuccess bool) error {
	if isSuccess {
		block.SuccessCount = block.SuccessCount + 1
	} else {
		block.FailureCount = block.FailureCount + 1
		if err := b.appDb.WithContext(ctx).Model(block).Association("FailureRanges").Append(failureRange); err != nil {
			b.logger.Error().Err(err).Str("domain", "block").Msg("failed to update block")
			return err
		}
	}

	return b.appDb.WithContext(ctx).Save(block).Error
}

func (b *BlockRepo) FindAllByBatchStatID(ctx context.Context, batchStatId string, offset, limit int) ([]*domain.Block, int64, error) {
	var blocks []*domain.Block
	var total int64

	tx := b.appDb.WithContext(ctx).Model(&domain.Block{})

	if err := tx.Where("batch_stat_id = ?", batchStatId).Count(&total).Error; err != nil {
		b.logger.Error().Err(err).Str("domain", "block").Str("batch_stat_id", batchStatId).Msg("failed to count block for specific batch")
		return nil, 0, err
	}

	if err := tx.Order("created_at DESC").Find(&blocks).Where("batch_stat_id = ?", batchStatId).Offset(offset).Limit(limit).Error; err != nil {
		b.logger.Error().Err(err).Str("domain", "block").Str("batch_stat_id", batchStatId).Msg("failed to fetch block for specific batch")
		return nil, 0, err
	}

	return blocks, total, nil
}

func (b *BlockRepo) FindById(ctx context.Context, uid string) (*domain.Block, error) {
	block := domain.Block{}

	if err := b.appDb.WithContext(ctx).First(&block, "id = ?", uid).Preload("FailureRanges").Error; err != nil {
		b.logger.Error().Err(err).Str("domain", "block").Str("id", uid).Msg("failed to fetch block for specific batch")
		return nil, err
	}

	return &block, nil
}
