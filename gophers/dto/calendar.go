package dto

import "context"

type CalendarService interface {
	ListEvents(ctx context.Context, calendarID string) ([]Event, error)
}

// Event represents a calendar event.
type Event struct {
	ID      string
	Summary string
	Start   string
	End     string
}
