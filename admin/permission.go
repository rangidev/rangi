package admin

import (
	"net/http"
	"time"
)

const (
	SessionCookieName     = "rangiAdminSession"
	SessionCookieDuration = 24 * time.Hour
)

func EnsurePermission(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: check permissions
		loggedIn := false
		sessionCookie, err := r.Cookie(SessionCookieName)
		if err == nil {
			if sessionCookie.Value == "true" {
				loggedIn = true
			}
		}
		if !loggedIn && r.URL.Path != LoginPath {
			// Not logged in and not trying to log in
			http.Redirect(w, r, LoginPath, http.StatusFound)
			return
		} else if loggedIn && r.URL.Path == LoginPath {
			// Logged in and trying to log in
			http.Redirect(w, r, DashboardPath, http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
