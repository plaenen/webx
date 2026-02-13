package webx

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
)

const (
	sessionCookieName = "webx_session"
	csrfSessionKey    = "csrf_token"
)

// SessionMiddleware reads or creates a session cookie, then populates
// WebXContext with the session ID and CSRF token from the store.
func SessionMiddleware(store SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionID, isNew := sessionIDFromRequest(r)

			if isNew {
				http.SetCookie(w, &http.Cookie{
					Name:     sessionCookieName,
					Value:    sessionID,
					Path:     "/",
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
				})
			}

			token, err := store.Get(sessionID, csrfSessionKey)
			if err != nil {
				http.Error(w, fmt.Sprintf("session store error: %v", err), http.StatusInternalServerError)
				return
			}

			if token == "" {
				token, err = randomHex(16)
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to generate CSRF token: %v", err), http.StatusInternalServerError)
					return
				}
				if err := store.Set(sessionID, csrfSessionKey, token); err != nil {
					http.Error(w, fmt.Sprintf("session store error: %v", err), http.StatusInternalServerError)
					return
				}
			}

			wctx := FromContext(r.Context())
			wctx.SessionID = sessionID
			wctx.CSRFToken = token

			next.ServeHTTP(w, r.WithContext(wctx.WithContext(r.Context())))
		})
	}
}

// sessionIDFromRequest returns the session ID from the cookie, or generates a
// new one. The bool indicates whether the ID is new.
func sessionIDFromRequest(r *http.Request) (string, bool) {
	if c, err := r.Cookie(sessionCookieName); err == nil && c.Value != "" {
		return c.Value, false
	}
	id, err := randomHex(16)
	if err != nil {
		// Extremely unlikely; fall back to a zero-value ID that will still work.
		return "0000000000000000", true
	}
	return id, true
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("reading random bytes: %w", err)
	}
	return hex.EncodeToString(b), nil
}
