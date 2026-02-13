package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestRangePage_Renders(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/range", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	ranges := page.Locator("main input.range")
	count, err := ranges.Count()
	if err != nil {
		t.Fatalf("count ranges: %v", err)
	}
	if count == 0 {
		t.Error("no range elements found")
	}
}

func TestRangePage_VariantStyles(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/range", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	variants := []string{
		"range-primary",
		"range-secondary",
		"range-accent",
		"range-info",
		"range-success",
		"range-warning",
		"range-error",
	}
	for _, v := range variants {
		loc := page.Locator("main input.range." + v)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", v, err)
		}
		if count == 0 {
			t.Errorf("no %s range found", v)
		}
	}
}

func TestRangePage_SizeVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/range", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	sizes := []string{"range-xs", "range-sm", "range-md", "range-lg", "range-xl"}
	for _, s := range sizes {
		loc := page.Locator("main input.range." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s range found", s)
		}
	}
}

func TestRangePage_WithSteps(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/range", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	stepped := page.Locator("main input.range[step='25']")
	count, err := stepped.Count()
	if err != nil {
		t.Fatalf("count stepped range: %v", err)
	}
	if count == 0 {
		t.Error("no range with step attribute found")
	}
}
