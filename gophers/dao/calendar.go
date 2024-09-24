package dao

import "time"

type Calendar struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"not null"`
	User       User   `gorm:"constraint:OnDelete:CASCADE;"`
	CalendarID string `gorm:"not null"`
	Name       string
	Events     []Event `gorm:"foreignKey:CalendarID;constraint:OnDelete:CASCADE;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
