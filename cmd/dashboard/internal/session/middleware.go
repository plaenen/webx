package session

import "net/http"

const sessionCookieName = "dashboard_session"

// Middleware reads the session ID from the cookie, looks up the UserSession
// from the store, and places it on the request context.
func Middleware(store SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionCookieName)
			if err != nil || cookie.Value == "" {
				// No session cookie — continue with nil session on context.
				next.ServeHTTP(w, r)
				return
			}

			userSession, err := store.Get(cookie.Value)
			if err != nil {
				http.Error(w, "session store error", http.StatusInternalServerError)
				return
			}

			if userSession == nil {
				// Session ID in cookie but not in store (expired/deleted).
				// Clear the stale cookie.
				http.SetCookie(w, &http.Cookie{
					Name:     sessionCookieName,
					Value:    "",
					Path:     "/",
					MaxAge:   -1,
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
				})
				next.ServeHTTP(w, r)
				return
			}

			ctx := WithUserSession(r.Context(), userSession)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// SetSessionCookie writes the session cookie to the response.
func SetSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearSessionCookie removes the session cookie from the response.
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
