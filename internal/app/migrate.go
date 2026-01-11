package app

import (
	"github.com/willbrid/api-gateway-sql/internal/domain"
	"github.com/willbrid/api-gateway-sql/pkg/logger"

	"gorm.io/gorm"
)

func MigrateAppDatabase(db *gorm.DB, iLogger logger.ILogger) {
	err := db.AutoMigrate(&domain.BatchStat{}, &domain.Block{}, &domain.FailureRange{})
	if err != nil {
		iLogger.Error("error during migrations: %s", err.Error())
	}
}
