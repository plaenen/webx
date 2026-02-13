package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestFileInputPage_Renders(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/file-input", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	inputs := page.Locator("main input.file-input")
	count, err := inputs.Count()
	if err != nil {
		t.Fatalf("count file-input: %v", err)
	}
	if count == 0 {
		t.Error("no file-input elements found")
	}
}

func TestFileInputPage_VariantStyles(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/file-input", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	variants := []string{
		"file-input-primary",
		"file-input-secondary",
		"file-input-accent",
		"file-input-info",
		"file-input-success",
		"file-input-warning",
		"file-input-error",
	}
	for _, v := range variants {
		loc := page.Locator("main input.file-input." + v)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", v, err)
		}
		if count == 0 {
			t.Errorf("no %s file-input found", v)
		}
	}
}

func TestFileInputPage_SizeVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/file-input", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	sizes := []string{"file-input-xs", "file-input-sm", "file-input-md", "file-input-lg", "file-input-xl"}
	for _, s := range sizes {
		loc := page.Locator("main input.file-input." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s file-input found", s)
		}
	}
}

func TestFileInputPage_Disabled(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/file-input", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	disabled := page.Locator("main input.file-input[disabled]")
	count, err := disabled.Count()
	if err != nil {
		t.Fatalf("count disabled: %v", err)
	}
	if count == 0 {
		t.Error("no disabled file-input found")
	}
}
