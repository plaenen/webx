package markdown

import (
	"bytes"
	"log"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// md is the configured goldmark instance.
// SECURITY: Raw HTML is NOT allowed - all HTML tags in markdown input
// are escaped to prevent XSS attacks. Use only markdown syntax.
var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM, // GitHub Flavored Markdown
		extension.Typographer,
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithXHTML(),
		// NOTE: html.WithUnsafe() is intentionally NOT used here.
		// Raw HTML in markdown is escaped to prevent XSS attacks.
	),
)

// Render converts markdown text to HTML.
// Raw HTML tags in the input are escaped for security.
func Render(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// MustRender converts markdown to HTML, returning empty string on error.
// Errors are logged for debugging but do not cause panics.
// Raw HTML tags in the input are escaped for security.
func MustRender(markdown string) string {
	result, err := Render(markdown)
	if err != nil {
		log.Printf("[markdown] render error: %v", err)
		return ""
	}
	return result
}
