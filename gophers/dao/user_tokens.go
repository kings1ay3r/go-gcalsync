package dao

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"time"
)

// TODO: Encrypt refresh token

// UserToken ...
type UserToken struct {
	AccountID    string    `gorm:"primaryKey;type:text"` // Account ID is part of the primary key
	UserID       int       `gorm:"primaryKey"`           // User ID is also part of the primary key
	AccessToken  string    `gorm:"type:text"`
	RefreshToken string    `gorm:"type:text"`
	Expiry       time.Time `gorm:"type:timestamp"`
}

// SaveUserTokens ...
func (d *dao) SaveUserTokens(ctx context.Context, userID int, accountID string, accessToken string, refreshToken string, expiry time.Time) error {

	var userTokens UserToken
	err := d.DB.Where("user_id = ? AND account_id = ?", userID, accountID).First(&userTokens).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			userTokens = UserToken{
				UserID:       userID,
				AccountID:    accountID,
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
				Expiry:       expiry,
			}
			return d.DB.Create(&userTokens).Error
		}
		return err
	}

	userTokens.AccessToken = accessToken
	userTokens.RefreshToken = refreshToken
	userTokens.Expiry = expiry

	return d.DB.Save(&userTokens).Error
}

// GetUserTokens ...
func (d *dao) GetUserTokens(ctx context.Context, userID int, accountID string) (*oauth2.Token, error) {
	var userTokens UserToken
	err := d.DB.Where("user_id = ? AND account_id = ?", userID, accountID).First(&userTokens).Error
	if err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken:  userTokens.AccessToken,
		RefreshToken: userTokens.RefreshToken,
		Expiry:       userTokens.Expiry,
	}, nil
}
