package internal

import "net/http"

// cacheControl wraps a handler to set Cache-Control: no-cache so browsers
// revalidate via ETags on each request.
func cacheControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		h.ServeHTTP(w, r)
	})
}
