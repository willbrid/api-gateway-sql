package app

import (
	"github.com/willbrid/api-gateway-sql/internal/domain"
	"github.com/willbrid/api-gateway-sql/pkg/logger"

	"gorm.io/gorm"
)

func MigrateAppDatabase(db *gorm.DB) {
	err := db.AutoMigrate(&domain.BatchStat{}, &domain.Block{}, &domain.FailureRange{})
	if err != nil {
		logger.Error("error during migrations: %s", err.Error())
	}
}
