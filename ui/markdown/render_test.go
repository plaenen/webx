package markdown

import (
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:     "heading",
			input:    "# Hello",
			contains: []string{"<h1", "Hello", "</h1>"},
		},
		{
			name:     "bold",
			input:    "**bold**",
			contains: []string{"<strong>", "bold", "</strong>"},
		},
		{
			name:     "italic",
			input:    "*italic*",
			contains: []string{"<em>", "italic", "</em>"},
		},
		{
			name:     "link",
			input:    "[example](https://example.com)",
			contains: []string{"<a", "href=\"https://example.com\"", "example", "</a>"},
		},
		{
			name:     "code",
			input:    "`code`",
			contains: []string{"<code>", "code", "</code>"},
		},
		{
			name:     "list",
			input:    "- item1\n- item2",
			contains: []string{"<ul>", "<li>", "item1", "item2", "</li>", "</ul>"},
		},
		{
			name:     "gfm strikethrough",
			input:    "~~deleted~~",
			contains: []string{"<del>", "deleted", "</del>"},
		},
		{
			name:     "gfm table",
			input:    "| A | B |\n|---|---|\n| 1 | 2 |",
			contains: []string{"<table>", "<th>", "A", "B", "<td>", "1", "2", "</table>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := Render(tt.input)
			if err != nil {
				t.Fatalf("Render(%q) error = %v", tt.input, err)
			}
			for _, s := range tt.contains {
				if !strings.Contains(html, s) {
					t.Errorf("Render(%q) = %q, want to contain %q", tt.input, html, s)
				}
			}
		})
	}
}

func TestMustRender(t *testing.T) {
	html := MustRender("# Hello")
	if !strings.Contains(html, "Hello") {
		t.Errorf("MustRender() = %q, want to contain 'Hello'", html)
	}

	html = MustRender("")
	if html != "" {
		t.Errorf("MustRender('') = %q, want empty string", html)
	}
}

func TestRender_SecurityHTMLEscaping(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		shouldNot  []string
		shouldHave []string
	}{
		{
			name:      "script tag injection",
			input:     "<script>alert('xss')</script>",
			shouldNot: []string{"<script>", "</script>", "alert("},
		},
		{
			name:      "onclick handler",
			input:     `<div onclick="alert('xss')">click me</div>`,
			shouldNot: []string{"onclick=", "alert("},
		},
		{
			name:      "img onerror",
			input:     `<img src="x" onerror="alert('xss')">`,
			shouldNot: []string{"onerror=", "alert("},
		},
		{
			name:      "iframe injection",
			input:     `<iframe src="https://evil.com"></iframe>`,
			shouldNot: []string{"<iframe", "</iframe>"},
		},
		{
			name:      "style tag",
			input:     `<style>body { display: none; }</style>`,
			shouldNot: []string{"<style>", "</style>"},
		},
		{
			name:      "link tag",
			input:     `<link rel="stylesheet" href="https://evil.com/steal.css">`,
			shouldNot: []string{"<link", "rel="},
		},
		{
			name:      "object tag",
			input:     `<object data="https://evil.com/exploit.swf"></object>`,
			shouldNot: []string{"<object", "</object>"},
		},
		{
			name:      "svg with script",
			input:     `<svg onload="alert('xss')"></svg>`,
			shouldNot: []string{"<svg", "onload="},
		},
		{
			name:      "mixed markdown and HTML",
			input:     "# Title\n\n<script>alert('xss')</script>\n\n**bold**",
			shouldNot: []string{"<script>", "</script>"},
			shouldHave: []string{"<h1", "Title", "<strong>", "bold"},
		},
		{
			name:      "HTML in code block is safe",
			input:     "```\n<script>alert('xss')</script>\n```",
			shouldNot: []string{},
			shouldHave: []string{"<code>", "&lt;script&gt;"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := Render(tt.input)
			if err != nil {
				t.Fatalf("Render(%q) error = %v", tt.input, err)
			}

			for _, s := range tt.shouldNot {
				if strings.Contains(html, s) {
					t.Errorf("SECURITY: Render(%q) = %q, should NOT contain %q", tt.input, html, s)
				}
			}

			for _, s := range tt.shouldHave {
				if !strings.Contains(html, s) {
					t.Errorf("Render(%q) = %q, should contain %q", tt.input, html, s)
				}
			}
		})
	}
}

func TestRender_SafeMarkdownFeatures(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:     "autolink URL",
			input:    "Visit https://example.com for more",
			contains: []string{"<a", "href=\"https://example.com\""},
		},
		{
			name:     "blockquote",
			input:    "> This is a quote",
			contains: []string{"<blockquote>", "This is a quote"},
		},
		{
			name:     "horizontal rule",
			input:    "---",
			contains: []string{"<hr"},
		},
		{
			name:     "nested list",
			input:    "- item1\n  - nested\n- item2",
			contains: []string{"<ul>", "<li>", "item1", "nested", "item2"},
		},
		{
			name:     "ordered list",
			input:    "1. first\n2. second",
			contains: []string{"<ol>", "<li>", "first", "second"},
		},
		{
			name:     "task list",
			input:    "- [ ] todo\n- [x] done",
			contains: []string{"<input", "type=\"checkbox\"", "todo", "done"},
		},
		{
			name:     "fenced code block with language",
			input:    "```go\nfunc main() {}\n```",
			contains: []string{"<code", "func main()"},
		},
		{
			name:     "image",
			input:    "![alt text](https://example.com/image.png)",
			contains: []string{"<img", "src=\"https://example.com/image.png\"", "alt=\"alt text\""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := Render(tt.input)
			if err != nil {
				t.Fatalf("Render(%q) error = %v", tt.input, err)
			}
			for _, s := range tt.contains {
				if !strings.Contains(html, s) {
					t.Errorf("Render(%q) = %q, want to contain %q", tt.input, html, s)
				}
			}
		})
	}
}
