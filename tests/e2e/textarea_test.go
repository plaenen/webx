package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestTextareaPage_TextareasRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/textarea", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	textareas := page.Locator("main textarea.textarea")
	count, err := textareas.Count()
	if err != nil {
		t.Fatalf("count textareas: %v", err)
	}
	if count == 0 {
		t.Error("no textarea components found on textarea page")
	}
}

func TestTextareaPage_ColorVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/textarea", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	colors := []string{
		"textarea-primary", "textarea-secondary", "textarea-accent", "textarea-neutral",
		"textarea-success", "textarea-warning", "textarea-info", "textarea-error",
	}
	for _, c := range colors {
		loc := page.Locator("main textarea." + c)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", c, err)
		}
		if count == 0 {
			t.Errorf("no %s textarea found", c)
		}
	}
}

func TestTextareaPage_SizeVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/textarea", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	sizes := []string{"textarea-xs", "textarea-sm", "textarea-md", "textarea-lg", "textarea-xl"}
	for _, s := range sizes {
		loc := page.Locator("main textarea." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s textarea found", s)
		}
	}
}

func TestTextareaPage_GhostVariant(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/textarea", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	loc := page.Locator("main textarea.textarea-ghost")
	count, err := loc.Count()
	if err != nil {
		t.Fatalf("count ghost: %v", err)
	}
	if count == 0 {
		t.Error("no ghost textarea found")
	}
}

func TestTextareaPage_DisabledPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/textarea", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	loc := page.Locator("main textarea.textarea[disabled]")
	count, err := loc.Count()
	if err != nil {
		t.Fatalf("count disabled: %v", err)
	}
	if count == 0 {
		t.Error("no disabled textarea found")
	}
}
