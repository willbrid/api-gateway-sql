package app

import (
	"api-gateway-sql/internal/domain"
	"api-gateway-sql/pkg/logger"

	"gorm.io/gorm"
)

func MigrateAppDatabase(db *gorm.DB) {
	err := db.AutoMigrate(&domain.BatchStat{}, &domain.Block{}, &domain.FailureRange{})
	if err != nil {
		logger.LogError("error during migrations: %s", err.Error())
	}
}
