package repository

import (
	"github.com/willbrid/api-gateway-sql/internal/domain"
	"github.com/willbrid/api-gateway-sql/pkg/uuid"

	"context"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type BatchStatRepo struct {
	appDb  *gorm.DB
	logger zerolog.Logger
}

func NewBatchStatRepo(appDb *gorm.DB, logger zerolog.Logger) *BatchStatRepo {
	return &BatchStatRepo{
		appDb:  appDb,
		logger: logger.With().Str("layer", "repository").Str("component", "batchstatrepo").Logger(),
	}
}

func (d *BatchStatRepo) Create(ctx context.Context, targetName string) (*domain.BatchStat, error) {
	uid := uuid.GenerateUID()

	batchStat := domain.BatchStat{
		ID:         uid,
		TargetName: targetName,
		Completed:  false,
	}

	if err := d.appDb.WithContext(ctx).Create(&batchStat).Error; err != nil {
		d.logger.Error().Err(err).Msg("failed to create batchStat")
		return nil, err
	}

	return &batchStat, nil
}

func (d *BatchStatRepo) UpdateLastCompleted(ctx context.Context, batchStat *domain.BatchStat) error {
	batchStat.Completed = true

	if err := d.appDb.WithContext(ctx).Save(&batchStat).Error; err != nil {
		d.logger.Error().Err(err).Msg("failed to update batchStat")
		return err
	}

	return nil
}

func (d *BatchStatRepo) AddBlockToBatchStat(ctx context.Context, bs *domain.BatchStat, block *domain.Block) (*domain.Block, error) {
	if err := d.appDb.WithContext(ctx).Model(bs).Association("Blocks").Append(block); err != nil {
		d.logger.Error().Err(err).Msg("failed to associate block to batchStat")
		return nil, err
	}

	if err := d.appDb.Save(bs).Error; err != nil {
		d.logger.Error().Err(err).Msg("failed to save batchStat after adding block")
		return nil, err
	}

	return block, nil
}

func (d *BatchStatRepo) FindById(ctx context.Context, uid string) (*domain.BatchStat, error) {
	batch := domain.BatchStat{}

	if err := d.appDb.WithContext(ctx).First(&batch, "id = ?", uid).Error; err != nil {
		d.logger.Error().Err(err).Str("batch_id", uid).Msg("failed to find batchStat by id")
		return nil, err
	}

	return &batch, nil
}

func (d *BatchStatRepo) CountUncompletedBatchStat(ctx context.Context) (int64, error) {
	var total int64

	err := d.appDb.WithContext(ctx).Model(&domain.BatchStat{}).Where("Completed = ?", false).Limit(1).Count(&total).Error
	if err != nil {
		d.logger.Error().Err(err).Msg("failed to count uncompleted batchStat")
		return 0, err
	}

	return total, nil
}

func (d *BatchStatRepo) FindAll(ctx context.Context, offset, limit int) ([]*domain.BatchStat, int64, error) {
	var batchStats []*domain.BatchStat
	var total int64

	tx := d.appDb.WithContext(ctx).Model(&domain.BatchStat{})

	if err := tx.Count(&total).Error; err != nil {
		d.logger.Error().Err(err).Msg("failed to count all batchStat")
		return nil, 0, err
	}

	if err := tx.Order("created_at DESC").Find(&batchStats).Offset(offset).Limit(limit).Error; err != nil {
		d.logger.Error().Err(err).Msg("failed to fetch batchStat")
		return nil, 0, err
	}

	return batchStats, total, nil
}
