package handlers

import "github.com/go-chi/chi/v5"

// Handlers wires all showcase SSE/API handlers.
type Handlers struct {
	validate  *validateHandlers
	parse     *parseHandlers
	form      *formHandlers
	upload    *uploadHandlers
	preview   *previewHandlers
}

func New() *Handlers {
	fileStore := newFileStore()
	return &Handlers{
		validate: newValidateHandlers(),
		parse:    newParseHandlers(),
		form:     newFormHandlers(),
		upload:   newUploadHandlers(fileStore),
		preview:  newPreviewHandlers(),
	}
}

// RegisterRoutes mounts all API handlers onto the given router.
func (h *Handlers) RegisterRoutes(r chi.Router) {
	h.validate.register(r)
	h.parse.register(r)
	h.form.register(r)
	h.upload.register(r)
	h.preview.register(r)
}
