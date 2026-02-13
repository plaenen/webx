package fileupload

import (
	"fmt"
	"net/http"
	"strings"

	webx "github.com/plaenen/webx"
	"github.com/plaenen/webx/utils"
	"github.com/starfederation/datastar-go/datastar"
)

// HandlerOption configures upload validation.
type HandlerOption func(*handlerConfig)

type handlerConfig struct {
	maxFileSize  int64    // per-file limit in bytes (default 10MB)
	allowedTypes []string // allowed MIME prefixes (empty = all)
	maxFiles     int      // max total files per component (default 10)
}

// WithMaxFileSize sets the maximum allowed size per file in bytes.
func WithMaxFileSize(bytes int64) HandlerOption {
	return func(c *handlerConfig) { c.maxFileSize = bytes }
}

// WithAllowedTypes restricts uploads to files whose Content-Type starts
// with one of the given prefixes (e.g. "image/", "application/pdf").
func WithAllowedTypes(types ...string) HandlerOption {
	return func(c *handlerConfig) { c.allowedTypes = types }
}

// WithMaxFiles sets the maximum number of files allowed per component.
func WithMaxFiles(n int) HandlerOption {
	return func(c *handlerConfig) { c.maxFiles = n }
}

func storeKey(sessionID, componentID string) string {
	return sessionID + ":" + componentID
}

// UploadHandler returns an http.HandlerFunc that accepts multipart file
// uploads and responds with an SSE patch of the updated file list.
//
// Mount at a dedicated POST path:
//
//	r.Post("/api/upload/files", fileupload.UploadHandler(store))
func UploadHandler(store *Store, opts ...HandlerOption) http.HandlerFunc {
	cfg := &handlerConfig{
		maxFileSize: 10 << 20, // 10MB
		maxFiles:    10,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		componentID := r.URL.Query().Get("id")
		if componentID == "" {
			http.Error(w, "missing id query parameter", http.StatusBadRequest)
			return
		}

		removeURL := r.URL.Query().Get("removeUrl")

		wctx := webx.FromContext(r.Context())
		key := storeKey(wctx.SessionID, componentID)

		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, fmt.Sprintf("parse form: %v", err), http.StatusBadRequest)
			return
		}

		existing := store.List(key)
		files := r.MultipartForm.File["files"]

		var errors []string
		for _, fh := range files {
			// Check max files
			if len(existing) >= cfg.maxFiles {
				errors = append(errors, fmt.Sprintf("maximum of %d files allowed", cfg.maxFiles))
				break
			}

			// Check file size
			if fh.Size > cfg.maxFileSize {
				errors = append(errors, fmt.Sprintf("%s exceeds maximum size", fh.Filename))
				continue
			}

			// Check MIME type
			ct := fh.Header.Get("Content-Type")
			if len(cfg.allowedTypes) > 0 {
				allowed := false
				for _, prefix := range cfg.allowedTypes {
					if strings.HasPrefix(ct, prefix) {
						allowed = true
						break
					}
				}
				if !allowed {
					errors = append(errors, fmt.Sprintf("%s: type %s not allowed", fh.Filename, ct))
					continue
				}
			}

			meta := FileMeta{
				ID:       utils.RandomID(),
				Name:     fh.Filename,
				Size:     fh.Size,
				MimeType: ct,
			}
			store.Add(key, meta)
			existing = append(existing, meta)
		}

		sse := datastar.NewSSE(w, r)
		if err := sse.PatchElementTempl(
			fileListItems(componentID, existing, removeURL),
			datastar.WithSelectorID(componentID+"-list"),
			datastar.WithModeInner(),
		); err != nil {
			return
		}

		if len(errors) > 0 {
			sse.PatchElements(
				fmt.Sprintf(`<p class="text-error text-sm">%s</p>`, strings.Join(errors, "; ")),
				datastar.WithSelectorID(componentID+"-errors"),
				datastar.WithModeInner(),
			)
		} else {
			sse.PatchElements(
				"",
				datastar.WithSelectorID(componentID+"-errors"),
				datastar.WithModeInner(),
			)
		}
	}
}

// RemoveHandler returns an http.HandlerFunc that removes a file from the
// store and responds with an SSE patch of the updated file list.
//
// Mount at a dedicated POST path:
//
//	r.Post("/api/upload/remove", fileupload.RemoveHandler(store))
func RemoveHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		componentID := r.URL.Query().Get("id")
		fileID := r.URL.Query().Get("fileId")
		removeURL := r.URL.Query().Get("removeUrl")
		if componentID == "" || fileID == "" {
			http.Error(w, "missing id or fileId query parameter", http.StatusBadRequest)
			return
		}

		wctx := webx.FromContext(r.Context())
		key := storeKey(wctx.SessionID, componentID)

		store.Remove(key, fileID)
		files := store.List(key)

		sse := datastar.NewSSE(w, r)
		sse.PatchElementTempl(
			fileListItems(componentID, files, removeURL),
			datastar.WithSelectorID(componentID+"-list"),
			datastar.WithModeInner(),
		)
	}
}
