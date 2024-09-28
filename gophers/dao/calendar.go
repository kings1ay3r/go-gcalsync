package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type Calendar struct {
	ID         uint   `gorm:"primaryKey"`
	CalendarID string `gorm:"not null;uniqueIndex:idx_user_calendar"` // Ensuring unique constraint with userID
	Name       string
	UserID     uint    `gorm:"not null"`                                           // Foreign key to User
	User       User    `gorm:"constraint:OnDelete:CASCADE;"`                       // Automatically delete calendar if a user is deleted
	Events     []Event `gorm:"foreignKey:CalendarID;constraint:OnDelete:CASCADE;"` // Automatically delete events if a calendar is deleted
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// FindCalendarByID checks if a calendar exists in the database by its CalendarID.
func FindCalendarByID(db *gorm.DB, calendarID string) (*Calendar, error) {
	var calendar Calendar
	err := db.Where("calendar_id = ?", calendarID).First(&calendar).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &calendar, nil
}

// InsertCalendar inserts a new calendar into the database.
func InsertCalendar(db *gorm.DB, calendarID, name string) (*Calendar, error) {
	calendar := &Calendar{
		CalendarID: calendarID,
		Name:       name,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err := db.Create(calendar).Error
	if err != nil {
		return nil, err
	}
	return calendar, nil
}

// FindCalendarByCalendarID ...
func (d *dao) FindCalendarByCalendarID(ctx context.Context, calendarID string) (*Calendar, error) {
	var calendar Calendar
	err := d.DB.Where("calendar_id = ?", calendarID).First(&calendar).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
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

	if err := d.DB.Where("user_id = ?", userID).Preload("Events").Find(&calendars).Error; err != nil {
		return nil, err
	}

	return calendars, nil
}
