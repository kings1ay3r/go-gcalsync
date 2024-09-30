package dao

import (
	"context"
	"fmt"
	"time"
)

// Watch ...
type Watch struct {
	ID         int    `gorm:"primary_key;auto_increment"`
	UserID     int    `gorm:"index:idx_user_calendar,unique"`
	CalendarID int    `gorm:"index:idx_user_calendar,unique"`
	ChannelID  string `gorm:"not null"`
	ResourceID string
	Expiry     time.Time `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Calendar Calendar `gorm:"constraint:OnDelete:CASCADE;"`
}

func (Watch) TableName() string {
	return "watches"
}

// SaveWatch upsert
func (d *dao) SaveWatch(ctx context.Context, watch *Watch) error {
	return d.DB.Save(watch).Error
}

// WatchWithDetails ...
type WatchesWithDetails struct {
	WatchID          int       `json:"watch_id"`
	UserID           int       `json:"user_id"`
	CalendarID       int       `json:"calendar_id"`
	GoogleCalendarID string    `json:"google_calendar_id"`
	ResourceID       string    `json:"resource_id"`
	AccountID        string    `json:"account_id"`
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	TokenExpiry      time.Time `json:"token_expiry"`
	Expiry           time.Time `json:"expiry"`
}

// FindExpiredWatches ...
func (d *dao) FindExpiringWatches(ctx context.Context) ([]WatchesWithDetails, error) {
	now := time.Now()
	tomorrowMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	var watches []WatchesWithDetails

	err := d.DB.Table("watches").
		Select("watches.id as watch_id, watches.user_id, watches.calendar_id, watches.resource_id, calendars.calendar_id as google_calendar_id, calendars.account_id, user_tokens.access_token, user_tokens.refresh_token, user_tokens.expiry as token_expiry, watches.expiry").
		Joins("JOIN calendars ON watches.calendar_id = calendars.id").
		Joins("JOIN user_tokens ON user_tokens.account_id = calendars.account_id AND user_tokens.user_id = calendars.user_id").
		Where("watches.expiry <= ?", tomorrowMidnight).
		Scan(&watches).
		Error

	// Process watchesWithDetails...

	if err != nil {
		return nil, fmt.Errorf("error fetching watches: %w", err)
	}

	return watches, err
}

// WatchExists ...
func (d *dao) WatchExists(ctx context.Context, userID uint, calendarID string) (bool, error) {
	var count int64
	err := d.DB.Model(&Watch{}).Where("user_id = ? AND calendar_id = ?", userID, calendarID).Count(&count).Error
	return count > 0, err
}
