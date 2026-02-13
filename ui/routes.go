package ui

import (
	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx/ui/calendar"
)

// RegisterRoutes registers all SSE/API handlers from UI component packages.
// Mount under your app's base path using chi's Route():
//
//	r.Route(basePath, func(r chi.Router) {
//	    ui.RegisterRoutes(r)
//	})
//
// For validator endpoints, mount them individually per validator:
//
//	r.Get("/api/validate/email", validator.Handler(emailValidator))
func RegisterRoutes(r chi.Router) {
	r.Get(calendar.NavigatePath, calendar.NavigateHandlerFromQuery())
}
