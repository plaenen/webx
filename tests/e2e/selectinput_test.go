package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestSelectPage_Renders(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/select", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	selects := page.Locator("main select.select")
	count, err := selects.Count()
	if err != nil {
		t.Fatalf("count selects: %v", err)
	}
	if count == 0 {
		t.Error("no select elements found")
	}
}

func TestSelectPage_VariantStyles(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/select", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	variants := []string{
		"select-primary",
		"select-secondary",
		"select-accent",
		"select-info",
		"select-success",
		"select-warning",
		"select-error",
		"select-ghost",
	}
	for _, v := range variants {
		loc := page.Locator("main select.select." + v)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", v, err)
		}
		if count == 0 {
			t.Errorf("no %s select found", v)
		}
	}
}

func TestSelectPage_SizeVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/select", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	sizes := []string{"select-xs", "select-sm", "select-md", "select-lg", "select-xl"}
	for _, s := range sizes {
		loc := page.Locator("main select.select." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s select found", s)
		}
	}
}

func TestSelectPage_Disabled(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/select", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	disabled := page.Locator("main select.select[disabled]")
	count, err := disabled.Count()
	if err != nil {
		t.Fatalf("count disabled: %v", err)
	}
	if count == 0 {
		t.Error("no disabled select found")
	}
}
