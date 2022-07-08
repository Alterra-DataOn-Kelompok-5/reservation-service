package domain

import (
	"time"
)

type Model struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
	// DeletedAt gorm.DeletedAt `gorm:"index"`
}
