package dao

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Calendar struct {
	ID         uint   `gorm:"primaryKey"`
	CalendarID string `gorm:"not null;uniqueIndex:idx_user_calendar"` // Ensuring unique constraint with userID
	Name       string
	UserID     uint    `gorm:"not null;uniqueIndex:idx_user_calendar"` // Foreign key to User
	User       User    `gorm:"constraint:OnDelete:CASCADE;"`
	Events     []Event `gorm:"foreignKey:CalendarID;constraint:OnDelete:CASCADE;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// FindCalendarByCalendarID ...
func (d *dao) FindCalendarByCalendarID(ctx context.Context, calendarID string) (*Calendar, error) {
	var calendar Calendar
	err := d.DB.Where("calendar_id = ?", calendarID).First(&calendar).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find calendar: %w", err)
	}
	return &calendar, nil
}

// InsertCalendar ...
func (d *dao) InsertCalendar(ctx context.Context, calendar Calendar) (*Calendar, error) {
	result := d.DB.Create(&calendar)
	if result.Error != nil {
		return nil, result.Error
	}
	return &calendar, nil
}

// GetUserCalendars ...
func (d *dao) GetUserCalendars(ctx context.Context, userID int) ([]Calendar, error) {
	var calendars []Calendar

	// TODO: Remove preloading. Implement a seperate paginated api for events. Control preload using a flag
	if err := d.DB.Where("user_id = ?", userID).Preload("Events").Find(&calendars).Error; err != nil {
		return nil, err
	}

	return calendars, nil
}
