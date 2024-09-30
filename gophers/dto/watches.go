package dto

import "time"

// WatchData holds user and calendar-related information
type WatchData struct {
	ID               int
	UserID           int
	CalendarID       int
	Expiry           time.Time
	ChannelID        string
	ResourceID       string
	AccountID        string
	GoogleCalendarID string
}

type RenewWatchesResponse struct {
	Succesful, Failed []*WatchData
}
