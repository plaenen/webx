package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestMarkdownPage_DisplayRendersHTML(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/markdown", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	display := page.Locator("#md-display")

	// Verify rendered markdown contains expected HTML elements
	h2 := display.Locator("h2")
	count, err := h2.Count()
	if err != nil {
		t.Fatalf("count h2: %v", err)
	}
	if count < 1 {
		t.Errorf("expected at least 1 <h2> in markdown display, got %d", count)
	}

	strong := display.Locator("strong")
	count, err = strong.Count()
	if err != nil {
		t.Fatalf("count strong: %v", err)
	}
	if count < 1 {
		t.Errorf("expected at least 1 <strong> in markdown display, got %d", count)
	}

	ul := display.Locator("ul")
	count, err = ul.Count()
	if err != nil {
		t.Fatalf("count ul: %v", err)
	}
	if count < 1 {
		t.Errorf("expected at least 1 <ul> in markdown display, got %d", count)
	}
}

func TestMarkdownPage_InputRendersTextarea(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/markdown", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	textarea := page.Locator("#md-editor")
	visible, err := textarea.IsVisible()
	if err != nil {
		t.Fatalf("is visible: %v", err)
	}
	if !visible {
		t.Errorf("expected textarea #md-editor to be visible")
	}

	placeholder, err := textarea.GetAttribute("placeholder")
	if err != nil {
		t.Fatalf("get placeholder: %v", err)
	}
	if placeholder != "Write markdown here..." {
		t.Errorf("expected placeholder 'Write markdown here...', got %q", placeholder)
	}
}

func TestMarkdownPage_PreviewTabShowsRendered(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/markdown", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Type markdown into the editor (use PressSequentially to fire input events)
	textarea := page.Locator("#md-editor")
	if err := textarea.Click(); err != nil {
		t.Fatalf("click textarea: %v", err)
	}
	if err := textarea.PressSequentially("# Test Heading", pw.LocatorPressSequentiallyOptions{
		Delay: pw.Float(20),
	}); err != nil {
		t.Fatalf("type: %v", err)
	}

	// Click Preview tab
	previewTab := page.Locator("#md-editor-container [role='tablist'] button:has-text('Preview')")
	if err := previewTab.Click(); err != nil {
		t.Fatalf("click preview tab: %v", err)
	}

	// Wait for preview to contain rendered content
	preview := page.Locator("#md-editor-preview")
	if err := preview.WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateVisible,
		Timeout: pw.Float(5000),
	}); err != nil {
		t.Fatalf("wait for preview visible: %v", err)
	}

	// Check that the preview contains rendered HTML
	h1 := preview.Locator("h1")
	if err := h1.WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateVisible,
		Timeout: pw.Float(5000),
	}); err != nil {
		t.Fatalf("wait for h1 in preview: %v", err)
	}

	text, err := h1.InnerText()
	if err != nil {
		t.Fatalf("h1 inner text: %v", err)
	}
	if text != "Test Heading" {
		t.Errorf("expected h1 'Test Heading', got %q", text)
	}
}

func TestMarkdownPage_WriteTabShowsTextarea(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/markdown", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Click Preview tab first
	previewTab := page.Locator("#md-editor-container [role='tablist'] button:has-text('Preview')")
	if err := previewTab.Click(); err != nil {
		t.Fatalf("click preview tab: %v", err)
	}

	// Verify textarea is hidden
	textarea := page.Locator("#md-editor")
	if err := textarea.WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateHidden,
		Timeout: pw.Float(3000),
	}); err != nil {
		t.Fatalf("wait for textarea hidden: %v", err)
	}

	// Click Write tab
	writeTab := page.Locator("#md-editor-container [role='tablist'] button:has-text('Write')")
	if err := writeTab.Click(); err != nil {
		t.Fatalf("click write tab: %v", err)
	}

	// Verify textarea is visible again
	if err := textarea.WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateVisible,
		Timeout: pw.Float(3000),
	}); err != nil {
		t.Fatalf("wait for textarea visible after Write click: %v", err)
	}

	visible, err := textarea.IsVisible()
	if err != nil {
		t.Fatalf("is visible: %v", err)
	}
	if !visible {
		t.Errorf("expected textarea to be visible after clicking Write tab")
	}
}
