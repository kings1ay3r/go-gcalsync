package dao

import "time"

type User struct {
	ID           uint       `gorm:"primaryKey"`
	Email        string     `gorm:"unique;not null"`
	PasswordHash string     `gorm:"not null"`
	Calendars    []Calendar `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"` // One-to-many relation (a user can have many calendars)
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
