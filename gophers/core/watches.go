package core

import (
	"context"
	"gcalsync/gophers/clients/logger"
	"gcalsync/gophers/dao"
	"gcalsync/gophers/dto"
	"golang.org/x/oauth2"
)

// RenewExpiringWatches ...
func (c *calendarClient) RenewExpiringWatches(ctx context.Context) {
	log := logger.GetInstance()
	expiringWatches, err := c.dao.FindExpiringWatches(ctx)
	if err != nil {
		// TODO: Implement channel logging
		log.Error(ctx, "unable to fetch pending watches")
	}
	var watchesToRenew []*dto.WatchData
	tokens := map[string]*oauth2.Token{}
	for _, watch := range expiringWatches {
		if watch.AccountID == "" {
			continue
		}
		token := &oauth2.Token{
			AccessToken:  watch.AccessToken,
			Expiry:       watch.Expiry,
			RefreshToken: watch.RefreshToken,
		}
		tokens[watch.AccountID] = token
		watchesToRenew = append(watchesToRenew, &dto.WatchData{
			ID:               watch.WatchID,
			UserID:           watch.UserID,
			CalendarID:       watch.CalendarID,
			Expiry:           watch.Expiry,
			AccountID:        watch.AccountID,
			GoogleCalendarID: watch.GoogleCalendarID,
			ResourceID:       watch.ResourceID,
		})

	}
	resp, errors := c.googleCalClient.RenewWatches(ctx, watchesToRenew, tokens)
	for _, err := range errors {
		log.Error(ctx, err.Error())
	}
	// TODO: Think about parallelizing / batchifying watches renewals, fault accommodation
	for _, rw := range resp.Succesfull {
		err := c.dao.SaveWatch(ctx, &dao.Watch{
			ID:         rw.ID,
			UserID:     rw.UserID,
			CalendarID: rw.CalendarID,
			ChannelID:  rw.ChannelID,
			ResourceID: rw.ResourceID,
			Expiry:     rw.Expiry,
		})
		if err != nil {
			log.Error(ctx, "unable to save watch #%d: %w", rw.ID, err)
		}
	}
	log.Info(ctx, "Renewed %d watch(es), failed to renew %d watch(es)", len(resp.Succesfull), len(resp.Failed))

}
