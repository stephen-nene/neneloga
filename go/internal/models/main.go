package models

import (
	"time"
)

// ServerStatus enum-like type
type ServerStatus string

const (
	StatusActive  ServerStatus = "active"
	StatusDeleted ServerStatus = "deleted"
	StatusOffline ServerStatus = "offline"
)

func (s ServerStatus) IsValid() bool {
	switch s {
	case StatusActive, StatusDeleted, StatusOffline:
		return true
	}
	return false
}

// AllStatuses returns all valid server statuses
func AllStatuses() []ServerStatus {
	return []ServerStatus{StatusActive, StatusDeleted, StatusOffline}
}

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"uniqueIndex;not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"default:user;not null"` // admin, user
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	Servers []Server `gorm:"foreignKey:UserID"`
}

type Log struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    *uint     `gorm:"index"`      // nullable
	ServerID  uint      `gorm:"not null;index"`
	Level     string    `gorm:"size:50;not null"`
	Message   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

type Server struct {
	ID        uint         `gorm:"primaryKey"`
	UserID    uint         `gorm:"not null;index"`
	Name      string       `gorm:"uniqueIndex;not null"`
	IPAddress string       `gorm:"uniqueIndex;not null"`
	Hostname  string       `gorm:"uniqueIndex;not null"`
	Os        string       `gorm:"not null"`
	Status    ServerStatus `gorm:"type:varchar(20);not null;default:active"`
	CreatedAt time.Time    `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time    `gorm:"not null" json:"updated_at"`

	// Relationships
	User User  `gorm:"foreignKey:UserID" json:"-"`
	Logs []Log `gorm:"foreignKey:ServerID" json:"logs,omitempty"`
}
