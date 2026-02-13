package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestLabelPage_Renders(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/label", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	labels := page.Locator("main span.label")
	count, err := labels.Count()
	if err != nil {
		t.Fatalf("count labels: %v", err)
	}
	if count == 0 {
		t.Error("no label elements found")
	}
}

func TestLabelPage_InsideInput(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/label", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	inputLabels := page.Locator("main label.input span.label")
	count, err := inputLabels.Count()
	if err != nil {
		t.Fatalf("count input labels: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 labels inside inputs, got %d", count)
	}
}

func TestLabelPage_InsideSelect(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/label", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	selectLabels := page.Locator("main label.select span.label")
	count, err := selectLabels.Count()
	if err != nil {
		t.Fatalf("count select labels: %v", err)
	}
	if count == 0 {
		t.Error("no labels inside select found")
	}
}

func TestLabelPage_FloatingLabel(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/label", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	floating := page.Locator("main label.floating-label")
	count, err := floating.Count()
	if err != nil {
		t.Fatalf("count floating-label: %v", err)
	}
	if count < 2 {
		t.Errorf("expected at least 2 floating labels, got %d", count)
	}
}
