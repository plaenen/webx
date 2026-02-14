package dashboard

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/plaenen/webx"
	"github.com/plaenen/webx/cmd/dashboard/internal/pages/dashboard/templates"
	"github.com/plaenen/webx/cmd/dashboard/internal/session"
	"github.com/plaenen/webx/idgen"
	"github.com/starfederation/datastar-go/datastar"
)

// SelfHostNewPageHandler renders the self-hosted provider creation form (GET).
func SelfHostNewPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		templates.SelfHostNew().Render(r.Context(), w)
	}
}

type selfHostSignals struct {
	SelfHost struct {
		Name string `json:"name"`
	} `json:"self_host"`
}

// SelfHostSubmitHandler creates a self-hosted provider with a pairing code (POST SSE).
func SelfHostSubmitHandler(
	providers webx.ProviderStore,
	selfHostType webx.ProviderType,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		us := session.GetUserSession(r.Context())
		if us == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var signals selfHostSignals
		if err := datastar.ReadSignals(r, &signals); err != nil {
			writeSelfHostError(w, r, "Failed to read form data")
			return
		}

		name := strings.TrimSpace(signals.SelfHost.Name)
		if name == "" {
			writeSelfHostError(w, r, "Name is required.")
			return
		}

		// Generate a 6-character pairing code formatted as XXX-XXX.
		code, err := generatePairingCode()
		if err != nil {
			writeSelfHostError(w, r, "Failed to generate pairing code.")
			return
		}

		prov := webx.Provider{
			ID:          idgen.Generate("prov"),
			Name:        name,
			Type:        selfHostType,
			OwnerID:     us.UserId,
			PairingCode: code,
			CreatedAt:   time.Now(),
		}

		if err := providers.Create(prov); err != nil {
			writeSelfHostError(w, r, "Failed to create provider.")
			return
		}

		sse := datastar.NewSSE(w, r)
		sse.MarshalAndPatchSignals(map[string]any{
			"self_host": map[string]any{
				"submitting":   false,
				"pairing_code": code,
				"provider_id":  prov.ID,
			},
		})
	}
}

// generatePairingCode creates a short code like "A1B-2C3".
func generatePairingCode() (string, error) {
	raw, err := idgen.Token(3)
	if err != nil {
		return "", fmt.Errorf("generating pairing code: %w", err)
	}
	upper := strings.ToUpper(raw)
	if len(upper) < 6 {
		return upper, nil
	}
	return upper[:3] + "-" + upper[3:6], nil
}

func writeSelfHostError(w http.ResponseWriter, r *http.Request, msg string) {
	sse := datastar.NewSSE(w, r)
	sse.MarshalAndPatchSignals(map[string]any{
		"self_host": map[string]any{
			"submitting": false,
			"error":      msg,
		},
	})
}
