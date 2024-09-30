package handlers

import (
	"fmt"
	"gcalsync/gophers/clients/logger"
	"net/http"
)

func (h *handler) NotifyHandler(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	resourceID := r.Header.Get("X-Goog-Resource-ID")
	resourceState := r.Header.Get("X-Goog-Resource-State")
	channelID := r.Header.Get("X-Goog-Channel-ID")
	expiry := r.Header.Get("X-Goog-Channel-Expiration")

	if resourceID == "" {
		logger.GetInstance().Info(nil, "resourceID not present.")
		return false, nil
	}

	logger.GetInstance().Info(nil, fmt.Sprintf("Received webhook: Resource ID: %s, State: %s, Channel ID: %s, Expiry: %s", resourceID, resourceState, channelID, expiry))

	if err := h.core.ProcessWebhook(r.Context(), resourceID); err != nil {
		logger.GetInstance().Error(nil, fmt.Sprintf("error processing webhook: %s", err.Error()))
	}

	return true, nil
}
