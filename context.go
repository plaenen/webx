// Provides context for webx templates
package webx

import (
	"context"
	"fmt"
)

type WebXContext struct {
	ShowDatastarInspector bool
	DatastarPro           bool
	CSRFToken             string
	DevMode               bool
	SessionID             string
	BasePath              string // prefix for all SSE handler routes (e.g. "/showcase")
}

func NewContext(ctx context.Context) *WebXContext {
	return &WebXContext{
		ShowDatastarInspector: false,
	}
}

type ctxKey struct{}

func FromContext(ctx context.Context) *WebXContext {
	if wctx, ok := ctx.Value(ctxKey{}).(*WebXContext); ok {
		return wctx
	}
	return NewContext(ctx)
}

func (wctx *WebXContext) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, wctx)
}

// APIPath returns the full path for a component SSE handler by prepending
// the BasePath. For example, with BasePath "/showcase" and path "/api/calendar/navigate",
// it returns "/showcase/api/calendar/navigate".
func (wctx *WebXContext) APIPath(path string) string {
	return wctx.BasePath + path
}

// Post returns a Datastar expression that performs a POST request to the given URL.
func Post(url string) string {
	return fmt.Sprintf("@post('%s')", url)
}
