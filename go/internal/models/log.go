package models

import (
	"time"
)

type Log struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"` // Foreign key
	Level     string    `gorm:"size:50;not null"`
	Message   string    `gorm:"type:text;not null"`
	CreatedAt time.Time
}
