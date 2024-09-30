package core

import (
	"context"
	"fmt"
	"gcalsync/gophers/clients/logger"
	"gcalsync/gophers/dao"
)

func (c *calendarClient) ProcessWebhook(ctx context.Context, resourceID string) error {
	log := logger.GetInstance()
	events, userID, err := c.googleCalClient.FetchEventsFromResource(ctx, resourceID)
	if err != nil {
		return fmt.Errorf("unable to retrieve calendar list: %w", err)
	}
	dbCalendar := &dao.CalendarData{Events: nil}
	for _, event := range events {
		dbEvent, err := mapEvent(event)

		if err != nil {
			log.Error(ctx, "failed to map event: %v", err)
		}
		dbCalendar.Events = append(dbCalendar.Events, dbEvent)
	}

	log.Info(ctx, "syncing %d events", len(events))

	// TODO: Track goroutines using ctx
	//#################################################################################
	//g, ctx := errgroup.WithContext(ctx)
	//g.Go(func() error {
	//	return c.insertInBackground(ctx, userID, []*dao.CalendarData{dbCalendar})
	//})
	//return g.Wait()
	//#################################################################################

	go c.insertInBackground(ctx, userID, []*dao.CalendarData{dbCalendar})

	return nil

}
