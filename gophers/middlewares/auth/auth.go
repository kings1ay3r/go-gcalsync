package auth

import (
	"context"
	"errors"
	"net/http"
)

const ContextUserKey = "currentUser"
const ContextUserIDKey = "currentUserID"

type Session struct {
	ID    int
	Email string
}

// Dummy User Model
var users = map[string]Session{
	"u1": {
		ID:    1,
		Email: "jane@example.com",
	},
	"u2": {
		ID:    2,
		Email: "john@example.com",
	},
}

// GetUserIDFromContext ...
func GetUserIDFromContext(ctx context.Context) (Session, error) {
	currUser, ok := ctx.Value(ContextUserKey).(Session)
	if !ok {
		return Session{}, errors.New("missing user context")
	}
	return currUser, nil
}

// BearerAuth ...
func BearerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		/***************************************************/
		/***************************************************/
		// FIXME: Bypass Auth. Need to Remove
		/***************************************************/

		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			userID = r.URL.Query().Get("state")
		}
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ContextUserIDKey, userID)
		user, ok := users[userID]
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx = context.WithValue(ctx, ContextUserKey, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
		return
		/***************************************************/
		/***************************************************/

		//authHeader := r.Header.Get("Authorization")
		//if !strings.HasPrefix(authHeader, "Bearer ") {
		//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//	return
		//}
		//
		//token := strings.TrimPrefix(authHeader, "Bearer ")
		//
		//user, ok := users[token]
		//if !ok {
		//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//	return
		//}
		//
		//ctx = context.WithValue(r.Context(), ContextUserKey, user)
		//r = r.WithContext(ctx)
		//next.ServeHTTP(w, r)
	})
}
