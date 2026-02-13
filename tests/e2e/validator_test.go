package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestValidatorPage_RendersInputs(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/validator", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	inputs := page.Locator("input.input")
	count, err := inputs.Count()
	if err != nil {
		t.Fatalf("count inputs: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 validator inputs, got %d", count)
	}
}

func TestValidatorPage_RendersHints(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/validator", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	hints := page.Locator("div.text-error")
	count, err := hints.Count()
	if err != nil {
		t.Fatalf("count hints: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 error hint elements, got %d", count)
	}
}

func TestValidatorPage_ShowsErrorOnInvalidEmail(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/validator", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Click input and type an invalid value character by character
	input := page.Locator("#email-demo")
	if err := input.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}
	if err := input.PressSequentially("bad", pw.LocatorPressSequentiallyOptions{
		Delay: pw.Float(50),
	}); err != nil {
		t.Fatalf("type: %v", err)
	}

	// Wait for the hint to become visible (debounced SSE call)
	hint := page.Locator("#email-demo-hint")
	if err := hint.WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateVisible,
		Timeout: pw.Float(5000),
	}); err != nil {
		t.Fatalf("wait for hint visible: %v", err)
	}

	visible, err := hint.IsVisible()
	if err != nil {
		t.Fatalf("is visible: %v", err)
	}
	if !visible {
		t.Errorf("expected error hint to be visible for invalid email")
	}
}

func TestValidatorPage_HidesErrorOnValidEmail(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/validator", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	input := page.Locator("#email-demo")

	// First type invalid to show error
	if err := input.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}
	if err := input.PressSequentially("bad", pw.LocatorPressSequentiallyOptions{
		Delay: pw.Float(50),
	}); err != nil {
		t.Fatalf("type bad: %v", err)
	}

	hint := page.Locator("#email-demo-hint")
	if err := hint.WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateVisible,
		Timeout: pw.Float(5000),
	}); err != nil {
		t.Fatalf("wait for hint visible: %v", err)
	}

	// Clear and type a valid email
	if err := input.Fill(""); err != nil {
		t.Fatalf("clear: %v", err)
	}
	if err := input.PressSequentially("test@example.com", pw.LocatorPressSequentiallyOptions{
		Delay: pw.Float(20),
	}); err != nil {
		t.Fatalf("type valid: %v", err)
	}

	// Wait for the hint to become hidden
	if err := hint.WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateHidden,
		Timeout: pw.Float(5000),
	}); err != nil {
		t.Fatalf("wait for hint hidden: %v", err)
	}

	visible, err := hint.IsVisible()
	if err != nil {
		t.Fatalf("is visible: %v", err)
	}
	if visible {
		t.Errorf("expected error hint to be hidden for valid email")
	}
}
