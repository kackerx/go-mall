package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type DemoOrder struct {
	ID        int64                 `gorm:"column:id;primary_key" json:"id"`
	UserID    int64                 `gorm:"column:user_id;not null;default:0" json:"user_id"`
	Amount    string                `gorm:"column:amout;not null;default:''" json:"amount"`
	Code      string                `gorm:"column:code;not null;default:'';unique_index" json:"code"`
	State     int8                  `gorm:"column:state;not null;default:0" json:"state"`
	PaidAt    time.Time             `gorm:"column:paid_at;default:'1970-01-01 00:00:00'" json:"paid_at"`
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag" json:"is_del"`
	CreatedAt time.Time             `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at" json:"updated_at"`
}
