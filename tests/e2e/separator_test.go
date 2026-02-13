package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestSeparatorPage_RendersBasicDivider(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/separator", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	dividers := page.Locator("div.divider")
	count, err := dividers.Count()
	if err != nil {
		t.Fatalf("count dividers: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 dividers, got %d", count)
	}
}

func TestSeparatorPage_ColorVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/separator", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	for _, variant := range []string{
		"divider-neutral",
		"divider-primary",
		"divider-secondary",
		"divider-accent",
		"divider-info",
		"divider-success",
		"divider-warning",
		"divider-error",
	} {
		loc := page.Locator("div.divider." + variant)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", variant, err)
		}
		if count < 1 {
			t.Errorf("expected at least 1 element with class %s, got %d", variant, count)
		}
	}
}

func TestSeparatorPage_HorizontalDirection(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/separator", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	horizontal := page.Locator("div.divider.divider-horizontal")
	count, err := horizontal.Count()
	if err != nil {
		t.Fatalf("count horizontal: %v", err)
	}
	if count < 1 {
		t.Errorf("expected at least 1 horizontal divider, got %d", count)
	}
}

func TestSeparatorPage_TextPositions(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/separator", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	start := page.Locator("div.divider.divider-start")
	startCount, err := start.Count()
	if err != nil {
		t.Fatalf("count start: %v", err)
	}
	if startCount < 1 {
		t.Errorf("expected at least 1 divider-start, got %d", startCount)
	}

	end := page.Locator("div.divider.divider-end")
	endCount, err := end.Count()
	if err != nil {
		t.Fatalf("count end: %v", err)
	}
	if endCount < 1 {
		t.Errorf("expected at least 1 divider-end, got %d", endCount)
	}
}
