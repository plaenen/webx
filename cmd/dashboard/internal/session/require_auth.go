package session

import "net/http"

// RequireAuth redirects unauthenticated requests to /login.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if us := GetUserSession(r.Context()); us != nil && us.IsLoggedIn() {
			next.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}
