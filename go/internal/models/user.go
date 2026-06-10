package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"uniqueIndex;not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"` // `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// A user can have many logs. This sets up the one-to-many relationship.
	Logs      []Log     `gorm:"foreignKey:UserID"`
}
