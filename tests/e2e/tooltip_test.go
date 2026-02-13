package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestTooltipPage_TooltipRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/tooltip", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main .tooltip")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count tooltip elements: %v", err)
	}
	if count == 0 {
		t.Error("no tooltip elements found on tooltip page")
	}
}

func TestTooltipPage_PositionVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/tooltip", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	positions := []string{"tooltip-top", "tooltip-bottom", "tooltip-left", "tooltip-right"}
	for _, pos := range positions {
		loc := page.Locator("main .tooltip." + pos)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", pos, err)
		}
		if count == 0 {
			t.Errorf("no %s tooltip found", pos)
		}
	}
}

func TestTooltipPage_ColorVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/tooltip", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	colors := []string{"tooltip-primary", "tooltip-secondary", "tooltip-accent", "tooltip-info", "tooltip-success", "tooltip-warning", "tooltip-error"}
	for _, c := range colors {
		loc := page.Locator("main .tooltip." + c)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", c, err)
		}
		if count == 0 {
			t.Errorf("no %s tooltip found", c)
		}
	}
}

func TestTooltipPage_CustomContent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/tooltip", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	loc := page.Locator("main .tooltip .tooltip-content")
	count, err := loc.Count()
	if err != nil {
		t.Fatalf("count tooltip-content: %v", err)
	}
	if count == 0 {
		t.Error("no tooltip-content elements found")
	}
}
