package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestMoneyInputPage_RendersInputs(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/money-input", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	inputs := page.Locator("input.input")
	count, err := inputs.Count()
	if err != nil {
		t.Fatalf("count inputs: %v", err)
	}
	if count < 4 {
		t.Errorf("expected at least 4 money/decimal inputs, got %d", count)
	}
}

func TestMoneyInputPage_DecimalParsesShorthand(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/money-input", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	input := page.Locator("#decimal-basic")
	if err := input.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}
	if err := input.PressSequentially("5k", pw.LocatorPressSequentiallyOptions{
		Delay: pw.Float(50),
	}); err != nil {
		t.Fatalf("type: %v", err)
	}

	amount := page.Locator("#decimal-basic-amount")
	if err := amount.WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateVisible,
		Timeout: pw.Float(5000),
	}); err != nil {
		t.Fatalf("wait for amount visible: %v", err)
	}

	text, err := amount.InnerText()
	if err != nil {
		t.Fatalf("inner text: %v", err)
	}
	if text != "5,000.00" {
		t.Errorf("expected formatted amount '5,000.00', got %q", text)
	}
}

func TestMoneyInputPage_DecimalShowsError(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/money-input", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	input := page.Locator("#decimal-basic")
	if err := input.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}
	if err := input.PressSequentially("bad", pw.LocatorPressSequentiallyOptions{
		Delay: pw.Float(50),
	}); err != nil {
		t.Fatalf("type: %v", err)
	}

	hint := page.Locator("#decimal-basic-hint")
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
		t.Errorf("expected error hint to be visible for invalid input")
	}
}

func TestMoneyInputPage_MoneyParsesCurrency(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/money-input", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	input := page.Locator("#money-any")
	if err := input.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}
	if err := input.PressSequentially("USD 5k", pw.LocatorPressSequentiallyOptions{
		Delay: pw.Float(50),
	}); err != nil {
		t.Fatalf("type: %v", err)
	}

	result := page.Locator("#money-any-result")
	if err := result.WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateVisible,
		Timeout: pw.Float(5000),
	}); err != nil {
		t.Fatalf("wait for result visible: %v", err)
	}

	badge := result.Locator(".badge")
	badgeText, err := badge.InnerText()
	if err != nil {
		t.Fatalf("badge inner text: %v", err)
	}
	if badgeText != "USD" {
		t.Errorf("expected currency badge 'USD', got %q", badgeText)
	}
}
