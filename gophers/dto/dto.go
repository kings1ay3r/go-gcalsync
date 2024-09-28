package dto

import "time"

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
}
