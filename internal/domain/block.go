package domain

import (
	"github.com/willbrid/api-gateway-sql/config"
	"github.com/willbrid/api-gateway-sql/pkg/csvstream"
	"github.com/willbrid/api-gateway-sql/pkg/uuid"
)

type BlockDataInput struct {
	BSInput *BatchStat
	BLInput *csvstream.Block
	TGInput *config.Target
	DBInput *config.Database
}

type Block struct {
	ID            string         `json:"id" gorm:"primaryKey"`
	StartLine     int            `json:"start_line"`
	EndLine       int            `json:"end_line"`
	SuccessCount  int            `json:"success" gorm:"default:0"`
	FailureCount  int            `json:"failure" gorm:"default:0"`
	FailureRanges []FailureRange `json:"failure_ranges" gorm:"foreignKey:BlockID"`
	CreatedAt     int64          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     int64          `json:"updated_at" gorm:"autoUpdateTime"`
	BatchStatID   string
}

func NewBlock(startLine, endLine int) *Block {
	return &Block{ID: uuid.GenerateUID(), StartLine: startLine, EndLine: endLine}
}
