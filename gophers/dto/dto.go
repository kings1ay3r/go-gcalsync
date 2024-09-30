package dto

import (
	"golang.org/x/oauth2"
	"time"
)

// Event represents a calendar event.
type Event struct {
	ID      string
	Summary string
	Start   time.Time
	End     time.Time
}

type User struct {
	ID int
}

type Calendar struct {
	Name   string
	ID     string
	Events []Event
	Token  *oauth2.Token
}

type CalendarDetailsByResourceIDResponse struct {
	ResourceID       string    `gorm:"column:resource_id"`
	UserID           int       `gorm:"column:user_id"`
	CalendarID       int       `gorm:"column:calendar_id"`
	GoogleCalendarID string    `gorm:"column:google_calendar_id"`
	AccountID        string    `gorm:"column:account_id"`
	AccessToken      string    `gorm:"column:access_token"`
	RefreshToken     string    `gorm:"column:refresh_token"`
	Expiry           time.Time `gorm:"column:expiry"`
	LastToken        string
}
