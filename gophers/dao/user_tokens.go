package dao

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"time"
)

// UserTokens represents the tokens for a user.
type UserTokens struct {
	UserID       int       `gorm:"primaryKey"`
	AccessToken  string    `gorm:"type:text"`
	RefreshToken string    `gorm:"type:text"`
	Expiry       time.Time `gorm:"type:timestamp"`
}

func (d *dao) SaveUserTokens(ctx context.Context, userID int, accessToken string, refreshToken string, expiry time.Time) error {

	var userTokens UserTokens
	err := d.DB.Where("user_id = ?", userID).First(&userTokens).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			userTokens = UserTokens{
				UserID:       userID,
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
func (d *dao) GetUserTokens(ctx context.Context, userID int) (*oauth2.Token, error) {
	var userTokens UserTokens
	err := d.DB.Where("user_id = ?", userID).First(&userTokens).Error
	if err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken:  userTokens.AccessToken,
		RefreshToken: userTokens.RefreshToken,
		Expiry:       userTokens.Expiry,
	}, nil
}
