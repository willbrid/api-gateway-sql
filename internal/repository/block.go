package repository

import (
	"github.com/willbrid/api-gateway-sql/internal/domain"

	"context"

	"gorm.io/gorm"
)

type BlockRepo struct {
	appDb *gorm.DB
}

func NewBlockRepo(appDb *gorm.DB) *BlockRepo {
	return &BlockRepo{appDb}
}

func (b *BlockRepo) Update(ctx context.Context, block *domain.Block, failureRange *domain.FailureRange, isSuccess bool) error {
	if isSuccess {
		block.SuccessCount = block.SuccessCount + 1
	} else {
		block.FailureCount = block.FailureCount + 1
		if err := b.appDb.WithContext(ctx).Model(block).Association("FailureRanges").Append(failureRange); err != nil {
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
		return nil, 0, err
	}

	if err := tx.Order("created_at DESC").Find(&blocks).Where("batch_stat_id = ?", batchStatId).Offset(offset).Limit(limit).Error; err != nil {
		return nil, 0, err
	}

	return blocks, total, nil
}

func (b *BlockRepo) FindById(ctx context.Context, uid string) (*domain.Block, error) {
	block := domain.Block{}

	if err := b.appDb.WithContext(ctx).First(&block, "id = ?", uid).Preload("FailureRanges").Error; err != nil {
		return nil, err
	}

	return &block, nil
}
