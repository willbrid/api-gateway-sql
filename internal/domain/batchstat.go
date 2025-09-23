package domain

type BatchStat struct {
	ID         string  `json:"id" gorm:"primaryKey"`
	TargetName string  `json:"target"`
	Completed  bool    `json:"completed"`
	CreatedAt  int64   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  int64   `json:"updated_at" gorm:"autoUpdateTime"`
	Blocks     []Block `json:"blocks" gorm:"foreignKey:BatchStatID"`
}
