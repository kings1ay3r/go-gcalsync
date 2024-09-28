package dao

import (
	"time"
)

type Event struct {
	ID          uint     `gorm:"primaryKey"`
	CalendarID  uint     `gorm:"not null"`                     // Foreign key to Calendar
	Calendar    Calendar `gorm:"constraint:OnDelete:CASCADE;"` // Automatically delete events if a calendar is deleted
	EventID     string   `gorm:"not null"`
	Summary     string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
