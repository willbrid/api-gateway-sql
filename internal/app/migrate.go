package app

import (
	"fmt"

	"github.com/willbrid/api-gateway-sql/internal/domain"

	"gorm.io/gorm"
)

func MigrateAppDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(&domain.BatchStat{}, &domain.Block{}, &domain.FailureRange{})
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return nil
}
