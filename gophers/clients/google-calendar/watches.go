package googlecalendar

import (
	"context"
	"fmt"
	"gcalsync/gophers/dto"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
	"time"
)

const (
	GOOGLE_WATCH_LIFE                                = time.Hour * 24 * 28
	__GOOGLE_UNSUPPORTED_CALENDAR_FOR_WATCH_REASON__ = "pushNotSupportedForRequestedResource"
	__ONE_HUNDRED_YEARS__                            = (100 * 365 * 24 * time.Hour) + (25 * 24 * time.Hour)
)

func (g *googleCalendar) RenewWatches(ctx context.Context, watches []*dto.WatchData, tokens map[string]*oauth2.Token) (res dto.RenewWatchesResponse, errs []error) {
	for _, w := range watches {

		token, _ := tokens[w.AccountID]
		_w, err := g.renewWatch(ctx, w, token)
		if err != nil {
			errs = append(errs, fmt.Errorf("unable to renew watch #%d: %w", w.ID, err))
			res.Failed = append(res.Failed, _w)
			continue
		}
		res.Succesfull = append(res.Succesfull, _w)
	}
	return
}

func (g *googleCalendar) renewWatch(ctx context.Context, watch *dto.WatchData, token *oauth2.Token) (*dto.WatchData, error) {

	// getOAuthClientForUser will find / renew token if it is missing
	client, err := g.getOAuthClientForUser(ctx, token, watch.UserID, watch.AccountID)
	if err != nil {
		return nil, fmt.Errorf("error creating OAuth client: %w", err)
	}

	calendarService, err := g.calendarService.NewService(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("error creating calendar service: %w", err)
	}

	channelID := getChannelID(watch)

	newWatch := &calendar.Channel{
		Id:         channelID,
		Type:       "web_hook",
		Address:    g.webhookURL,
		Expiration: time.Now().Add(GOOGLE_WATCH_LIFE).UnixMilli(),
		Params: map[string]string{
			"token": g.webhookURL,
		},
	}

	watchResp, err := g.calendarService.Watch(ctx, calendarService, watch.GoogleCalendarID, newWatch)
	if err != nil {
		if gErr, ok := err.(*googleapi.Error); ok {
			if len(gErr.Errors) == 1 && gErr.Errors[0].Reason == __GOOGLE_UNSUPPORTED_CALENDAR_FOR_WATCH_REASON__ {
				return &dto.WatchData{
					ID:         watch.ID,
					UserID:     watch.UserID,
					CalendarID: watch.CalendarID,
					Expiry:     time.Now().Add(__ONE_HUNDRED_YEARS__),
					ResourceID: "0",
					AccountID:  watch.AccountID,
					ChannelID:  channelID,
				}, nil
			}
		}
		return nil, fmt.Errorf("error renewing watch: %w", err)
	}

	expiry := time.UnixMilli(watchResp.Expiration)

	return &dto.WatchData{
		ID:         watch.ID,
		UserID:     watch.UserID,
		CalendarID: watch.CalendarID,
		Expiry:     expiry,
		ResourceID: watchResp.ResourceId,
		AccountID:  watch.AccountID,
		ChannelID:  channelID,
	}, nil
}

func getChannelID(watch *dto.WatchData) string {
	now := time.Now()
	channelID := fmt.Sprintf("user-%d-calendar-%d-%d", watch.UserID, watch.CalendarID, now.UnixMilli())
	return channelID
}
