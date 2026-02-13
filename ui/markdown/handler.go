package markdown

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/starfederation/datastar-go/datastar"
)

type previewSignals struct {
	Value string `json:"value"`
	Mode  string `json:"mode"`
}

// PreviewHandler returns an http.HandlerFunc that renders markdown from
// the component's signals and patches the preview div via SSE.
//
// Mount at a dedicated POST path:
//
//	r.Post("/api/preview/markdown", markdown.PreviewHandler())
func PreviewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		componentID := r.URL.Query().Get("id")
		if componentID == "" {
			http.Error(w, "missing id query parameter", http.StatusBadRequest)
			return
		}

		sanitizedID := strings.ReplaceAll(componentID, "-", "_")
		wrapper := map[string]previewSignals{}
		if err := datastar.ReadSignals(r, &wrapper); err != nil {
			http.Error(w, fmt.Sprintf("read signals: %v", err), http.StatusBadRequest)
			return
		}

		store, ok := wrapper[sanitizedID]
		if !ok {
			http.Error(w, fmt.Sprintf("missing signals for %q", sanitizedID), http.StatusBadRequest)
			return
		}

		html, err := Render(store.Value)
		if err != nil {
			html = fmt.Sprintf(`<p class="text-error text-sm">Render error: %s</p>`, err.Error())
		}
		if store.Value == "" {
			html = `<p class="text-base-content/50 italic">Nothing to preview</p>`
		}

		sse := datastar.NewSSE(w, r)
		sse.PatchElements(
			html,
			datastar.WithSelectorID(componentID+"-preview"),
			datastar.WithModeInner(),
		)
	}
}
