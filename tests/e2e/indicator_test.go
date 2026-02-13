package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestIndicatorPage_RendersBadgeOnBox(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/indicator", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	indicators := page.Locator("div.indicator")
	count, err := indicators.Count()
	if err != nil {
		t.Fatalf("count indicators: %v", err)
	}
	if count < 1 {
		t.Errorf("expected at least 1 indicator, got %d", count)
	}
}

func TestIndicatorPage_RendersIndicatorItems(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/indicator", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	items := page.Locator("span.indicator-item")
	count, err := items.Count()
	if err != nil {
		t.Fatalf("count indicator-items: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 indicator-items, got %d", count)
	}
}

func TestIndicatorPage_AllPositions(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/indicator", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	for _, cls := range []string{
		"indicator-start",
		"indicator-center",
		"indicator-middle",
		"indicator-bottom",
	} {
		loc := page.Locator("span.indicator-item." + cls)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", cls, err)
		}
		if count < 1 {
			t.Errorf("expected at least 1 element with class %s, got %d", cls, count)
		}
	}
}

func TestIndicatorPage_BadgeOnButton(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/indicator", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// An indicator wrapping a button should exist
	btnInIndicator := page.Locator("div.indicator button.btn")
	count, err := btnInIndicator.Count()
	if err != nil {
		t.Fatalf("count button in indicator: %v", err)
	}
	if count < 1 {
		t.Errorf("expected at least 1 button inside an indicator, got %d", count)
	}
}
