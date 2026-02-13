package calendar

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/starfederation/datastar-go/datastar"
)

// NavigatePath is the standard handler path for calendar navigation.
// Mount it under your app's base path: basePath + NavigatePath.
const NavigatePath = "/api/calendar/navigate"

// navigateSignals is the signal shape sent by the client during navigation.
type navigateSignals struct {
	Year      int    `json:"year"`
	Month     int    `json:"month"`
	Direction int    `json:"direction"` // -1 = prev, +1 = next
	Selected  string `json:"selected"`
}

// NavigateHandler returns an http.HandlerFunc that handles SSE-based
// month navigation for a calendar component. The calendarID must match
// the ID used when rendering the Calendar component so that
// PatchElementTempl can morph the correct DOM node.
func NavigateHandler(calendarID string, mode Mode) http.HandlerFunc {
	return handleNavigate(calendarID, mode)
}

// NavigateHandlerFromQuery returns an http.HandlerFunc that reads the
// calendar ID and mode from query parameters "id" and "mode". This is
// useful when a single endpoint serves multiple calendar instances.
func NavigateHandlerFromQuery() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		calendarID := r.URL.Query().Get("id")
		if calendarID == "" {
			http.Error(w, "missing id query parameter", http.StatusBadRequest)
			return
		}

		mode := ModeSingle
		if r.URL.Query().Get("mode") == "range" {
			mode = ModeRange
		}

		handleNavigate(calendarID, mode)(w, r)
	}
}

func handleNavigate(calendarID string, mode Mode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Signals arrive namespaced: {"cal_id": {"year": ..., "month": ..., "direction": ...}}
		wrapper := map[string]navigateSignals{}
		if err := datastar.ReadSignals(r, &wrapper); err != nil {
			http.Error(w, fmt.Sprintf("read signals: %v", err), http.StatusBadRequest)
			return
		}

		// Extract the namespaced signals using sanitized ID (hyphens → underscores).
		sanitizedID := strings.ReplaceAll(calendarID, "-", "_")
		store, ok := wrapper[sanitizedID]
		if !ok {
			http.Error(w, fmt.Sprintf("missing signals for %q", sanitizedID), http.StatusBadRequest)
			return
		}

		// Compute new month/year.
		t := time.Date(store.Year, time.Month(store.Month), 1, 0, 0, 0, 0, time.UTC)
		t = t.AddDate(0, store.Direction, 0)
		newYear := t.Year()
		newMonth := t.Month()

		// Build the calendar props for re-rendering.
		props := Props{
			ID:       calendarID,
			Year:     newYear,
			Month:    newMonth,
			Selected: store.Selected,
			Mode:     mode,
		}

		// Create SSE writer and send the patched element + signals.
		sse := datastar.NewSSE(w, r)

		if err := sse.PatchElementTempl(Calendar(props)); err != nil {
			return
		}

		// Patch the signals so the client knows the new year/month.
		updatedSignals := map[string]any{
			sanitizedID: map[string]any{
				"year":  newYear,
				"month": int(newMonth),
			},
		}
		sse.MarshalAndPatchSignals(updatedSignals)
	}
}
