package model

import (
	"time"

	"gorm.io/gorm"
)

type Common struct {
	ID        uint            `json:"id"`
	CreatedAt time.Time       `json:"created_at" gorm:"type:datetime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"type:datetime"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" gorm:"type:datetime"`
}

func (c *Common) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	return
}

func (c *Common) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return
}
