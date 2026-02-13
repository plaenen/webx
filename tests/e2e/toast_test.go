package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestToastPage_ToastsRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/toast", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	toasts := page.Locator(".toast")
	count, err := toasts.Count()
	if err != nil {
		t.Fatalf("count toasts: %v", err)
	}
	if count == 0 {
		t.Error("no toast components found on toast page")
	}
}

func TestToastPage_PositionsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/toast", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	positions := []string{
		"toast-start",
		"toast-center",
		"toast-end",
		"toast-top",
		"toast-middle",
		"toast-bottom",
	}
	for _, p := range positions {
		loc := page.Locator("." + p)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", p, err)
		}
		if count == 0 {
			t.Errorf("no %s toast found", p)
		}
	}
}

func TestToastPage_ContainsAlerts(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/toast", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Each toast should contain at least one alert child.
	alertsInToasts := page.Locator(".toast .alert")
	count, err := alertsInToasts.Count()
	if err != nil {
		t.Fatalf("count alerts in toasts: %v", err)
	}
	if count == 0 {
		t.Error("no alert children found inside toast containers")
	}
}

func TestToastPage_StackedAlerts(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/toast", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The stacked demo has a single toast with multiple alerts.
	// Find any toast that has more than one alert.
	toasts := page.Locator(".toast")
	count, err := toasts.Count()
	if err != nil {
		t.Fatalf("count toasts: %v", err)
	}

	foundStacked := false
	for i := 0; i < count; i++ {
		alerts := toasts.Nth(i).Locator(".alert")
		alertCount, err := alerts.Count()
		if err != nil {
			t.Fatalf("count alerts in toast %d: %v", i, err)
		}
		if alertCount > 1 {
			foundStacked = true
			break
		}
	}
	if !foundStacked {
		t.Error("no stacked toast (with multiple alerts) found")
	}
}
