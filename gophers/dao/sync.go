package dao

import (
	"context"
	"errors"
	"gcalsync/gophers/clients/logger"
	"gorm.io/gorm"
	"time"
)

// SaveUserCalendarData saves the user’s calendar and associated events in a single transaction.
func (d *dao) SaveUserCalendarData(ctx context.Context, userID int, calendars []*CalendarData) error {
	log := logger.GetInstance()
	// TODO: Refactor to use GORM Transaction
	tx := d.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, calendarData := range calendars {
		var calendar Calendar
		// Check if the calendar already exists
		err := tx.Where("calendar_id = ? AND user_id = ?", calendarData.CalendarID, userID).First(&calendar).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return rollbackWithError(tx, ctx, err, "error searching calendar")
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Calendar does not exist, create a new one
			calendar = Calendar{
				CalendarID: calendarData.CalendarID,
				Name:       calendarData.Name,
				UserID:     uint(userID),
				AccountID:  calendarData.AccountID,
			}
			if err := tx.Create(&calendar).Error; err != nil {
				return rollbackWithError(tx, ctx, err, "error searching calendar")
			}
			log.Info(ctx, "save user calendar data: %v", calendar.ID)
		} else {
			log.Info(ctx, "Calendar Exists: %v", calendar.CalendarID)
		}

		for _, eventData := range calendarData.Events {
			var event Event
			err = tx.Where("calendar_id = ? AND event_id = ?", calendar.ID, eventData.EventID).First(&event).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return rollbackWithError(tx, ctx, err, "error searching calendar")
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				event = Event{
					CalendarID: calendar.ID,
					EventID:    eventData.EventID,
					Summary:    eventData.Name,
					StartTime:  eventData.StartTime,
					EndTime:    eventData.EndTime,
				}
				if err := tx.Create(&event).Error; err != nil {
					return rollbackWithError(tx, ctx, err, "error searching calendar")
				}
				log.Info(ctx, "Event created: %v", event.ID)
			} else {
				event.Summary = eventData.Name
				event.StartTime = eventData.StartTime
				event.EndTime = eventData.EndTime
				event.UpdatedAt = time.Now()
				if err := tx.Save(&event).Error; err != nil {
					return rollbackWithError(tx, ctx, err, "error searching calendar")
				}
				log.Info(ctx, "Event updated: %v", event.ID)
			}
		}

		// TODO: Dont save unsupported watches. eg; en.indian#holiday@group.v.calendar.google.com
		watch := &Watch{
			UserID:     userID,
			CalendarID: calendar.ID,
			Expiry:     time.Now(),
		}
		if err := tx.Save(&watch).Error; err != nil {
			return rollbackWithError(tx, ctx, err, "error inserting watch")
		}
		err = err
	}

	return tx.Commit().Error
}

func rollbackWithError(tx *gorm.DB, ctx context.Context, err error, message string) error {
	tx.Rollback()
	logger.GetInstance().Error(ctx, message+": %v", err)
	return err
}

// CalendarData holds the data for a calendar and its associated events.
type CalendarData struct {
	Name       string
	CalendarID string
	Events     []EventData
	AccountID  string
}

// EventData holds the data for an event.
type EventData struct {
	EventID   string
	Name      string
	StartTime time.Time
	EndTime   time.Time
}
