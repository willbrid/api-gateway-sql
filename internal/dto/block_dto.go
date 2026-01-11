package dto

import (
	"github.com/willbrid/api-gateway-sql/config"
	"github.com/willbrid/api-gateway-sql/internal/domain"
	"github.com/willbrid/api-gateway-sql/pkg/csvstream"
)

type BlockDataInput struct {
	BSInput *domain.BatchStat
	BLInput *csvstream.Block
	TGInput *config.Target
	DBInput *config.Database
}
