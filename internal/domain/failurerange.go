package domain

import "github.com/willbrid/api-gateway-sql/pkg/uuid"

type FailureRange struct {
	ID        string `json:"id" gorm:"primaryKey"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	BlockID   string
}

func NewFailureRange(startLine, endLine int) *FailureRange {
	return &FailureRange{ID: uuid.GenerateUID(), StartLine: startLine, EndLine: endLine}
}
