package middleware

import (
	"net/http"

	"backend/utils"
)

// RequireAdmin allows access only when the session_role cookie is admin.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_role")
		if err != nil || cookie.Value != "admin" {
			utils.Error(w, http.StatusUnauthorized, "admin access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}
