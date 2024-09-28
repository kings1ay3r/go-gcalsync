package handlers

import (
	"errors"
	"fmt"
	"gcalsync/gophers/middlewares/auth"
	"net/http"
)

// CallbackHandler handles the OAuth2 callback from Google
func (h *handler) CallbackHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	code := r.URL.Query().Get("code")
	ctx := r.Context()
	err := h.core.InsertCalendars(ctx, code)
	if err != nil {
		return nil, errors.New("unable to insert calendars")
	}

	currUserID := ctx.Value(auth.ContextUserIDKey)

	redirectUrl := fmt.Sprintf("/index?user_id=%s", currUserID.(string))

	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
	return nil, nil
}
