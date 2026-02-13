// Package ds provides typed helpers for Datastar HTML attributes.
//
// Datastar's parameterized plugins use a colon separator (e.g. data-on:click),
// NOT a hyphen (data-on-click). A hyphen is silently ignored as an unknown
// plugin. This package makes that mistake impossible by construction.
package ds

import (
	"fmt"
	"strings"

	"github.com/a-h/templ"
)

// --- Parameterized attributes (colon syntax) ---

// On returns a data-on:<event> attribute.
//
//	ds.On("click", expr) → {"data-on:click": expr}
func On(event, expr string) templ.Attributes {
	return templ.Attributes{"data-on:" + event: expr}
}

// OnClick is shorthand for On("click", expr).
func OnClick(expr string) templ.Attributes {
	return On("click", expr)
}

// Bind returns a data-bind:<signal> attribute.
func Bind(signal string) templ.Attributes {
	return templ.Attributes{"data-bind:" + signal: ""}
}

// ClassToggle returns a data-class:<name> attribute (single class toggle).
func ClassToggle(name, expr string) templ.Attributes {
	return templ.Attributes{"data-class:" + name: expr}
}

// Attr returns a data-attr:<name> attribute.
func Attr(name, expr string) templ.Attributes {
	return templ.Attributes{"data-attr:" + name: expr}
}

// Style returns a data-style:<prop> attribute.
func Style(prop, expr string) templ.Attributes {
	return templ.Attributes{"data-style:" + prop: expr}
}

// Computed returns a data-computed:<name> attribute.
func Computed(name, expr string) templ.Attributes {
	return templ.Attributes{"data-computed:" + name: expr}
}

// Indicator returns a data-indicator:<name> attribute.
func Indicator(name string) templ.Attributes {
	return templ.Attributes{"data-indicator:" + name: ""}
}

// Ref returns a data-ref:<name> attribute.
func Ref(name string) templ.Attributes {
	return templ.Attributes{"data-ref:" + name: ""}
}

// --- Standalone attributes (no colon) ---

// Signals returns a data-signals attribute.
func Signals(value string) templ.Attributes {
	return templ.Attributes{"data-signals": value}
}

// Show returns a data-show attribute.
func Show(expr string) templ.Attributes {
	return templ.Attributes{"data-show": expr}
}

// Text returns a data-text attribute.
func Text(expr string) templ.Attributes {
	return templ.Attributes{"data-text": expr}
}

// Class returns a data-class attribute (object syntax).
func Class(value string) templ.Attributes {
	return templ.Attributes{"data-class": value}
}

// Init returns a data-init attribute.
func Init(expr string) templ.Attributes {
	return templ.Attributes{"data-init": expr}
}

// Effect returns a data-effect attribute.
func Effect(expr string) templ.Attributes {
	return templ.Attributes{"data-effect": expr}
}

// --- Backend action expressions ---
//
// These return JS expressions for use in data-on:* handlers.
// Mutating methods (POST, PUT, PATCH, DELETE) automatically include
// the CSRF token from <meta name="csrf-token">.
//
// By default, Datastar's built-in retry behavior is used (retry: 'auto',
// retryMaxCount: 10). Use WithRetries to customize, or the *Once
// convenience functions for single-shot requests.

// csrfJS is the JS expression that reads the CSRF token from the meta tag.
const csrfJS = `document.querySelector('meta[name=csrf-token]')?.content||''`

// ActionOption customizes a backend action expression (@get, @post, etc.).
type ActionOption func(*actionConfig)

type actionConfig struct {
	retries     *int   // nil = Datastar default; 0 = no retry; >0 = custom count
	contentType string // e.g. "form" for multipart/form-data
}

// WithRetries sets the maximum number of retry attempts.
// Use 0 to disable retries entirely.
func WithRetries(n int) ActionOption {
	return func(c *actionConfig) { c.retries = &n }
}

// WithContentType sets the content type for the action.
// Use "form" to send multipart/form-data (file uploads).
func WithContentType(ct string) ActionOption {
	return func(c *actionConfig) { c.contentType = ct }
}

// noRetry is a pre-built option that disables retries.
var noRetry = WithRetries(0)

// buildAction constructs a @method('url', {options}) expression.
func buildAction(method, url string, csrf bool, opts []ActionOption) string {
	cfg := &actionConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	var parts []string
	if csrf {
		parts = append(parts, fmt.Sprintf("headers: {'X-CSRF-Token': %s}", csrfJS))
	}
	if cfg.contentType != "" {
		parts = append(parts, fmt.Sprintf("contentType: '%s'", cfg.contentType))
	}
	if cfg.retries != nil {
		parts = append(parts, fmt.Sprintf("retryMaxCount: %d", *cfg.retries))
	}

	if len(parts) == 0 {
		return fmt.Sprintf("@%s('%s')", method, url)
	}
	return fmt.Sprintf("@%s('%s', {%s})", method, url, strings.Join(parts, ", "))
}

// Get returns a @get('url') expression.
//
//	ds.Get("/api/data")                // → @get('/api/data')
//	ds.Get("/api/data", ds.WithRetries(3)) // → @get('/api/data', {retryMaxCount: 3})
func Get(url string, opts ...ActionOption) string {
	return buildAction("get", url, false, opts)
}

// GetOnce returns a @get('url') expression without retries.
//
//	ds.GetOnce("/api/data") // → @get('/api/data', {retryMaxCount: 0})
func GetOnce(url string) string {
	return Get(url, noRetry)
}

// Post returns a @post('url') expression with the CSRF token header.
//
//	ds.Post("/api/submit")                // → @post('/api/submit', {headers: {…}})
//	ds.Post("/api/submit", ds.WithRetries(5)) // → @post('/api/submit', {headers: {…}, retryMaxCount: 5})
func Post(url string, opts ...ActionOption) string {
	return buildAction("post", url, true, opts)
}

// PostOnce returns a @post('url') expression with CSRF but without retries.
func PostOnce(url string) string {
	return Post(url, noRetry)
}

// Put returns a @put('url') expression with the CSRF token header.
func Put(url string, opts ...ActionOption) string {
	return buildAction("put", url, true, opts)
}

// PutOnce returns a @put('url') expression with CSRF but without retries.
func PutOnce(url string) string {
	return Put(url, noRetry)
}

// Patch returns a @patch('url') expression with the CSRF token header.
func Patch(url string, opts ...ActionOption) string {
	return buildAction("patch", url, true, opts)
}

// PatchOnce returns a @patch('url') expression with CSRF but without retries.
func PatchOnce(url string) string {
	return Patch(url, noRetry)
}

// Delete returns a @delete('url') expression with the CSRF token header.
func Delete(url string, opts ...ActionOption) string {
	return buildAction("delete", url, true, opts)
}

// DeleteOnce returns a @delete('url') expression with CSRF but without retries.
func DeleteOnce(url string) string {
	return Delete(url, noRetry)
}

// --- Combiner ---

// Merge combines multiple templ.Attributes into one.
// Later values overwrite earlier ones for the same key.
func Merge(attrs ...templ.Attributes) templ.Attributes {
	merged := templ.Attributes{}
	for _, attr := range attrs {
		for k, v := range attr {
			merged[k] = v
		}
	}
	return merged
}
