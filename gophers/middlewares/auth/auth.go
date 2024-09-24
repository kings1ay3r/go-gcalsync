package auth

import (
	"net/http"
)

func BearerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//authHeader := r.Header.Get("Authorization")
		//if !strings.HasPrefix(authHeader, "Bearer ") {
		//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//	return
		//}
		//
		//token := strings.TrimPrefix(authHeader, "Bearer ")
		//if token != os.Getenv("BEARER_TOKEN") {
		//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//	return
		//}

		next.ServeHTTP(w, r)
	})
}
