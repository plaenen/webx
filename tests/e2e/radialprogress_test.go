package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestRadialProgressPage_Render(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/radial-progress", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main .radial-progress")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count radial-progress elements: %v", err)
	}
	if count == 0 {
		t.Error("no radial-progress elements found")
	}
}

func TestRadialProgressPage_HasAriaAttributes(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/radial-progress", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	loc := page.Locator("main .radial-progress[role='progressbar'][aria-valuenow]")
	count, err := loc.Count()
	if err != nil {
		t.Fatalf("count accessible radial-progress: %v", err)
	}
	if count == 0 {
		t.Error("no radial-progress elements with role and aria-valuenow found")
	}
}

func TestRadialProgressPage_MultipleValues(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/radial-progress", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main .radial-progress")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count radial-progress: %v", err)
	}
	if count < 5 {
		t.Errorf("expected at least 5 radial-progress elements, got %d", count)
	}
}
