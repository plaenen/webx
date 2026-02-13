package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx/ui/fileupload"
)

type uploadHandlers struct {
	store *fileupload.Store
}

func newFileStore() *fileupload.Store {
	return fileupload.NewStore()
}

func newUploadHandlers(store *fileupload.Store) *uploadHandlers {
	return &uploadHandlers{store: store}
}

func (u *uploadHandlers) register(r chi.Router) {
	r.Post("/api/upload/files", fileupload.UploadHandler(u.store))
	r.Post("/api/upload/files-restricted", fileupload.UploadHandler(u.store,
		fileupload.WithAllowedTypes("image/"),
		fileupload.WithMaxFiles(3),
	))
	r.Post("/api/upload/remove", fileupload.RemoveHandler(u.store))
}
