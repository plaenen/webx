package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestMockupCodePage_Renders(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/mockup-code", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main div.mockup-code")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count mockup-code: %v", err)
	}
	if count == 0 {
		t.Error("no mockup-code containers found")
	}
}

func TestMockupCodePage_LinePrefix(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/mockup-code", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	prefixed := page.Locator("main .mockup-code pre[data-prefix]")
	count, err := prefixed.Count()
	if err != nil {
		t.Fatalf("count prefixed lines: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 prefixed lines, got %d", count)
	}
}

func TestMockupCodePage_CodeElements(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/mockup-code", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	codes := page.Locator("main .mockup-code pre code")
	count, err := codes.Count()
	if err != nil {
		t.Fatalf("count code elements: %v", err)
	}
	if count < 5 {
		t.Errorf("expected at least 5 code elements, got %d", count)
	}
}

func TestMockupCodePage_HighlightedLine(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/mockup-code", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	highlighted := page.Locator("main .mockup-code pre.bg-warning")
	count, err := highlighted.Count()
	if err != nil {
		t.Fatalf("count highlighted: %v", err)
	}
	if count == 0 {
		t.Error("no highlighted line found")
	}
}
