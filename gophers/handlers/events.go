package handlers

import (
	"net/http"
)

func (h *handler) ListEventsHandler(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()
	return h.core.GetMyCalendarEvents(ctx)
}

// ConnectHandler initiates the OAuth2 flow for Google Calendar access
func (h *handler) ConnectHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()
	url, err := h.core.GetAuthCodeURL(ctx)
	if err != nil {
		http.Error(w, "Failed to initiate OAuth flow", http.StatusInternalServerError)
		return nil, err
	}
	http.Redirect(w, r, url, http.StatusSeeOther)
	return nil, nil
}
