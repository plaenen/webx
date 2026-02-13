package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestMoneyPage_Renders(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/money", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	spans := page.Locator("main span.inline-flex")
	count, err := spans.Count()
	if err != nil {
		t.Fatalf("count money spans: %v", err)
	}
	if count == 0 {
		t.Error("no money elements found")
	}
}

func TestMoneyPage_FormattedValues(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/money", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Check that thousands separator is present
	content, err := page.Locator("main").TextContent()
	if err != nil {
		t.Fatalf("get text content: %v", err)
	}

	values := []string{"1,234.56", "1,000,000.00", "9,999.99"}
	for _, v := range values {
		if !contains(content, v) {
			t.Errorf("expected formatted value %q in page content", v)
		}
	}
}

func TestMoneyPage_CurrencySymbols(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/money", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	currencies := []string{"$", "EUR", "£"}
	content, err := page.Locator("main").TextContent()
	if err != nil {
		t.Fatalf("get text content: %v", err)
	}
	for _, c := range currencies {
		if !contains(content, c) {
			t.Errorf("expected currency %q in page content", c)
		}
	}
}

func TestMoneyPage_Precision(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/money", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	content, err := page.Locator("main").TextContent()
	if err != nil {
		t.Fatalf("get text content: %v", err)
	}

	// Zero precision: "1,235" (rounded)
	if !contains(content, "1,235") {
		t.Error("expected zero-precision value 1,235")
	}
	// Four precision: "1,234.5678"
	if !contains(content, "1,234.5678") {
		t.Error("expected four-precision value 1,234.5678")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
