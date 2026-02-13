package validator

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/starfederation/datastar-go/datastar"
)

// Result holds the outcome of a validation check.
type Result struct {
	Valid bool
	Error string
}

// ValidateFunc validates a string value and returns a Result.
type ValidateFunc func(value string) Result

// Handler returns an http.HandlerFunc for a single ValidateFunc.
// The component ID is passed as a query parameter "id".
//
// Mount each validator at its own path:
//
//	r.Get("/api/validate/email", validator.Handler(emailValidator))
//	r.Get("/api/validate/phone", validator.Handler(phoneValidator))
//
// The component references this path via the ValidateURL prop.
func Handler(fn ValidateFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		componentID := r.URL.Query().Get("id")
		if componentID == "" {
			http.Error(w, "missing id query parameter", http.StatusBadRequest)
			return
		}

		// Signals arrive namespaced: {"component_id": {"value": "...", ...}}
		sanitizedID := strings.ReplaceAll(componentID, "-", "_")
		wrapper := map[string]inputSignals{}
		if err := datastar.ReadSignals(r, &wrapper); err != nil {
			http.Error(w, fmt.Sprintf("read signals: %v", err), http.StatusBadRequest)
			return
		}

		store, ok := wrapper[sanitizedID]
		if !ok {
			http.Error(w, fmt.Sprintf("missing signals for %q", sanitizedID), http.StatusBadRequest)
			return
		}

		result := fn(store.Value)

		sse := datastar.NewSSE(w, r)
		sse.MarshalAndPatchSignals(map[string]any{
			sanitizedID: map[string]any{
				"valid": result.Valid,
				"error": result.Error,
			},
		})
	}
}

// inputSignals is the signal shape sent by the client.
type inputSignals struct {
	Value string `json:"value"`
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}
