package dao

import "time"

type User struct {
	ID           int    `gorm:"primaryKey"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Calendars  []Calendar  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"` // One-to-many relation (a user can have many calendars)
	UserTokens []UserToken `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Watches    []Watch     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}
