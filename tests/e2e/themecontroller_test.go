package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestThemeController_ToggleSwitchesTheme(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/theme-controller", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The toggle input should exist with the theme-controller class.
	toggle := page.Locator("input.toggle.theme-controller[value='dark']")
	if err := toggle.WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateVisible,
		Timeout: pw.Float(3000),
	}); err != nil {
		t.Fatalf("toggle not found: %v", err)
	}

	// Click the toggle to switch to dark theme.
	if err := toggle.Click(); err != nil {
		t.Fatalf("click toggle: %v", err)
	}

	// Verify the document root has data-theme="dark".
	theme, err := page.Evaluate("document.documentElement.getAttribute('data-theme')")
	if err != nil {
		t.Fatalf("get data-theme: %v", err)
	}
	if theme != "dark" {
		t.Errorf("expected data-theme='dark', got %q", theme)
	}

	// Click again to toggle back to default.
	if err := toggle.Click(); err != nil {
		t.Fatalf("click toggle again: %v", err)
	}

	theme, err = page.Evaluate("document.documentElement.getAttribute('data-theme')")
	if err != nil {
		t.Fatalf("get data-theme after untoggle: %v", err)
	}
	if theme != "default" {
		t.Errorf("expected data-theme='default' after untoggle, got %q", theme)
	}
}

func TestThemeController_RadioGroupSwitchesTheme(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/theme-controller", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Find the radio group container.
	radioGroup := page.Locator("[data-signals*='tc_radio']")
	if err := radioGroup.WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateVisible,
		Timeout: pw.Float(3000),
	}); err != nil {
		t.Fatalf("radio group not found: %v", err)
	}

	// Click the "retro" radio button.
	retroRadio := radioGroup.Locator("input.radio[value='retro']")
	if err := retroRadio.Click(); err != nil {
		t.Fatalf("click retro radio: %v", err)
	}

	// Verify theme changed to retro.
	theme, err := page.Evaluate("document.documentElement.getAttribute('data-theme')")
	if err != nil {
		t.Fatalf("get data-theme: %v", err)
	}
	if theme != "retro" {
		t.Errorf("expected data-theme='retro', got %q", theme)
	}

	// Click cyberpunk to switch again.
	cyberpunkRadio := radioGroup.Locator("input.radio[value='cyberpunk']")
	if err := cyberpunkRadio.Click(); err != nil {
		t.Fatalf("click cyberpunk radio: %v", err)
	}

	theme, err = page.Evaluate("document.documentElement.getAttribute('data-theme')")
	if err != nil {
		t.Fatalf("get data-theme after cyberpunk: %v", err)
	}
	if theme != "cyberpunk" {
		t.Errorf("expected data-theme='cyberpunk', got %q", theme)
	}
}

func TestThemeController_ButtonGroupSwitchesTheme(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/theme-controller", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Find the button group container (join class).
	buttonGroup := page.Locator("[data-signals*='tc_buttons']")
	if err := buttonGroup.WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateVisible,
		Timeout: pw.Float(3000),
	}); err != nil {
		t.Fatalf("button group not found: %v", err)
	}

	// Click the "synthwave" button.
	synthBtn := buttonGroup.Locator("input.btn[value='synthwave']")
	if err := synthBtn.Click(); err != nil {
		t.Fatalf("click synthwave button: %v", err)
	}

	// Verify theme changed.
	theme, err := page.Evaluate("document.documentElement.getAttribute('data-theme')")
	if err != nil {
		t.Fatalf("get data-theme: %v", err)
	}
	if theme != "synthwave" {
		t.Errorf("expected data-theme='synthwave', got %q", theme)
	}
}

func TestThemeController_HasExpectedInputs(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/theme-controller", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Toggle section should have exactly one checkbox input.
	toggleInputs := page.Locator("input.toggle.theme-controller")
	count, err := toggleInputs.Count()
	if err != nil {
		t.Fatalf("count toggle inputs: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 toggle input, got %d", count)
	}

	// Radio group should have 6 radio inputs (one per theme).
	radioInputs := page.Locator("[data-signals*='tc_radio'] input.radio.theme-controller")
	radioCount, err := radioInputs.Count()
	if err != nil {
		t.Fatalf("count radio inputs: %v", err)
	}
	if radioCount != 6 {
		t.Errorf("expected 6 radio inputs, got %d", radioCount)
	}

	// Button group should have 6 button inputs.
	btnInputs := page.Locator("[data-signals*='tc_buttons'] input.btn.theme-controller")
	btnCount, err := btnInputs.Count()
	if err != nil {
		t.Fatalf("count button inputs: %v", err)
	}
	if btnCount != 6 {
		t.Errorf("expected 6 button inputs, got %d", btnCount)
	}
}
