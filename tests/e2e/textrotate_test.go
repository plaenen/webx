package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestTextRotatePage_Render(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/text-rotate", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main .text-rotate")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count text-rotate: %v", err)
	}
	if count == 0 {
		t.Error("no text-rotate elements found")
	}
}

func TestTextRotatePage_MultipleInstances(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/text-rotate", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main .text-rotate")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count text-rotate: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 text-rotate elements, got %d", count)
	}
}

func TestTextRotatePage_HasTextItems(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/text-rotate", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	items := page.Locator("main .text-rotate > span > span")
	count, err := items.Count()
	if err != nil {
		t.Fatalf("count text items: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 rotating text items, got %d", count)
	}
}
