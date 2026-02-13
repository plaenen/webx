package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx/ui/markdown"
)

type previewHandlers struct{}

func newPreviewHandlers() *previewHandlers {
	return &previewHandlers{}
}

func (p *previewHandlers) register(r chi.Router) {
	r.Post("/api/preview/markdown", markdown.PreviewHandler())
}
