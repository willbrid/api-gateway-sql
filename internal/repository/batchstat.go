package repository

import (
	"api-gateway-sql/internal/domain"
	"api-gateway-sql/pkg/uuid"

	"context"

	"gorm.io/gorm"
)

type BatchStatRepo struct {
	appDb *gorm.DB
}

func NewBatchStatRepo(appDb *gorm.DB) *BatchStatRepo {
	return &BatchStatRepo{appDb}
}

func (d *BatchStatRepo) Create(ctx context.Context, targetName string) (string, error) {
	uid := uuid.GenerateUID()

	batchStat := domain.BatchStat{
		ID:         uid,
		TargetName: targetName,
		Completed:  false,
	}

	err := d.appDb.WithContext(ctx).Create(&batchStat).Error

	return uid, err
}

func (d *BatchStatRepo) UpdateLastCompleted(ctx context.Context, batchStat *domain.BatchStat) error {
	batchStat.Completed = true

	return d.appDb.WithContext(ctx).Save(&batchStat).Error
}

func (d *BatchStatRepo) AddBlockToBatchStat(ctx context.Context, bs *domain.BatchStat, block *domain.Block) (*domain.Block, error) {
	if err := d.appDb.WithContext(ctx).Model(bs).Association("Blocks").Append(block); err != nil {
		return nil, err
	}

	if err := d.appDb.Save(bs).Error; err != nil {
		return nil, err
	}

	return block, nil
}

func (d *BatchStatRepo) FindById(ctx context.Context, uid string) (*domain.BatchStat, error) {
	batch := domain.BatchStat{}

	if err := d.appDb.WithContext(ctx).First(&batch, uid).Error; err != nil {
		return nil, err
	}

	return &batch, nil
}

func (d *BatchStatRepo) FindAll(ctx context.Context, offset, limit int) ([]*domain.BatchStat, int64, error) {
	var batchStats []*domain.BatchStat
	var total int64

	tx := d.appDb.WithContext(ctx).Model(&domain.BatchStat{})

	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := tx.Find(&batchStats).Offset(offset).Limit(limit).Error; err != nil {
		return nil, 0, err
	}

	return batchStats, total, nil
}
