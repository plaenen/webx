package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestTogglePage_TogglesRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/toggle", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	toggles := page.Locator("main .toggle")
	count, err := toggles.Count()
	if err != nil {
		t.Fatalf("count toggles: %v", err)
	}
	if count == 0 {
		t.Error("no toggle components found on toggle page")
	}
}

func TestTogglePage_ColorVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/toggle", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	colors := []string{
		"toggle-primary", "toggle-secondary", "toggle-accent", "toggle-neutral",
		"toggle-success", "toggle-warning", "toggle-info", "toggle-error",
	}
	for _, c := range colors {
		loc := page.Locator("main .toggle." + c)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", c, err)
		}
		if count == 0 {
			t.Errorf("no %s toggle found", c)
		}
	}
}

func TestTogglePage_SizeVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/toggle", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	sizes := []string{"toggle-xs", "toggle-sm", "toggle-md", "toggle-lg", "toggle-xl"}
	for _, s := range sizes {
		loc := page.Locator("main .toggle." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s toggle found", s)
		}
	}
}

func TestTogglePage_DisabledPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/toggle", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	loc := page.Locator("main .toggle[disabled]")
	count, err := loc.Count()
	if err != nil {
		t.Fatalf("count disabled: %v", err)
	}
	if count == 0 {
		t.Error("no disabled toggle found")
	}
}
