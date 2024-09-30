package dao

import (
	"time"
)

type Event struct {
	ID          int    `gorm:"primaryKey"`
	CalendarID  int    `gorm:"not null"` // Foreign key to Calendar
	EventID     string `gorm:"not null"`
	Summary     string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	Calendar Calendar `gorm:"constraint:OnDelete:CASCADE;"`
}
