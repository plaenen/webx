package e2e_test

import (
	"fmt"
	"testing"
	"time"

	pw "github.com/playwright-community/playwright-go"
)

func TestCalendarAdvanced_NavigationRendersNewMonth(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/calendar-advanced", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Get the current month label from the navigation header.
	now := time.Now()
	currentLabel := fmt.Sprintf("%s %d", now.Month().String(), now.Year())

	navHeader := page.Locator("[data-signals*='nav_cal'] .font-semibold").First()
	text, err := navHeader.TextContent()
	if err != nil {
		t.Fatalf("get nav header text: %v", err)
	}
	if text != currentLabel {
		t.Errorf("initial month label = %q, want %q", text, currentLabel)
	}

	// Click Next button.
	nextBtn := page.Locator("[data-signals*='nav_cal'] button", pw.PageLocatorOptions{
		HasText: "Next",
	}).First()
	if err := nextBtn.Click(); err != nil {
		t.Fatalf("click next: %v", err)
	}

	// Wait for the month label to change.
	nextMonth := now.AddDate(0, 1, 0)
	expectedLabel := fmt.Sprintf("%s %d", nextMonth.Month().String(), nextMonth.Year())

	if err := page.Locator("[data-signals*='nav_cal'] .font-semibold", pw.PageLocatorOptions{
		HasText: expectedLabel,
	}).First().WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateVisible,
		Timeout: pw.Float(5000),
	}); err != nil {
		t.Fatalf("next month label %q not visible: %v", expectedLabel, err)
	}
}

func TestCalendarAdvanced_NavigationPrevNext(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/calendar-advanced", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	now := time.Now()
	currentLabel := fmt.Sprintf("%s %d", now.Month().String(), now.Year())

	// Click Prev.
	prevBtn := page.Locator("[data-signals*='nav_cal'] button", pw.PageLocatorOptions{
		HasText: "Prev",
	}).First()
	if err := prevBtn.Click(); err != nil {
		t.Fatalf("click prev: %v", err)
	}

	prevMonth := now.AddDate(0, -1, 0)
	prevLabel := fmt.Sprintf("%s %d", prevMonth.Month().String(), prevMonth.Year())

	if err := page.Locator("[data-signals*='nav_cal'] .font-semibold", pw.PageLocatorOptions{
		HasText: prevLabel,
	}).First().WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateVisible,
		Timeout: pw.Float(5000),
	}); err != nil {
		t.Fatalf("prev month label %q not visible: %v", prevLabel, err)
	}

	// Click Next to go back to original month.
	nextBtn := page.Locator("[data-signals*='nav_cal'] button", pw.PageLocatorOptions{
		HasText: "Next",
	}).First()
	if err := nextBtn.Click(); err != nil {
		t.Fatalf("click next: %v", err)
	}

	if err := page.Locator("[data-signals*='nav_cal'] .font-semibold", pw.PageLocatorOptions{
		HasText: currentLabel,
	}).First().WaitFor(pw.LocatorWaitForOptions{
		State:   pw.WaitForSelectorStateVisible,
		Timeout: pw.Float(5000),
	}); err != nil {
		t.Fatalf("return to current month %q not visible: %v", currentLabel, err)
	}
}

func TestCalendarAdvanced_RangeSelection(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/calendar-advanced", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Find the range calendar section.
	rangeCal := page.Locator("#range-cal")

	// Click day 10 (range start).
	day10 := rangeCal.Locator("button", pw.LocatorLocatorOptions{
		HasText: "10",
	}).First()
	if err := day10.Click(); err != nil {
		t.Fatalf("click day 10: %v", err)
	}

	// Click day 20 (range end).
	day20 := rangeCal.Locator("button", pw.LocatorLocatorOptions{
		HasText: "20",
	}).First()
	if err := day20.Click(); err != nil {
		t.Fatalf("click day 20: %v", err)
	}

	// Both start and end days should have btn-primary class.
	for _, dayText := range []string{"10", "20"} {
		btn := rangeCal.Locator("button.btn-primary", pw.LocatorLocatorOptions{
			HasText: dayText,
		}).First()
		if err := btn.WaitFor(pw.LocatorWaitForOptions{
			State:   pw.WaitForSelectorStateVisible,
			Timeout: pw.Float(3000),
		}); err != nil {
			t.Errorf("day %s not highlighted with btn-primary: %v", dayText, err)
		}
	}
}

func TestCalendarAdvanced_RangeHighlighting(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/calendar-advanced", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	rangeCal := page.Locator("#range-cal")

	// Click day 10 then day 15 to create a range.
	day10 := rangeCal.Locator("button", pw.LocatorLocatorOptions{
		HasText: "10",
	}).First()
	if err := day10.Click(); err != nil {
		t.Fatalf("click day 10: %v", err)
	}

	day15 := rangeCal.Locator("button", pw.LocatorLocatorOptions{
		HasText: "15",
	}).First()
	if err := day15.Click(); err != nil {
		t.Fatalf("click day 15: %v", err)
	}

	// Days 11-14 should have accent styling (btn-accent).
	accentBtns := rangeCal.Locator("button.btn-accent")
	count, err := accentBtns.Count()
	if err != nil {
		t.Fatalf("count accent buttons: %v", err)
	}
	if count < 1 {
		t.Error("expected in-between days to have btn-accent class, found none")
	}
}

func TestCalendarAdvanced_RangeStateIndicator(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/calendar-advanced", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The range badge should initially show "none".
	badge := page.Locator("[data-signals*='range_cal'] code.badge")
	text, err := badge.TextContent()
	if err != nil {
		t.Fatalf("badge text: %v", err)
	}
	if text != "none" {
		t.Errorf("initial range badge = %q, want %q", text, "none")
	}

	// Click a day in the range calendar.
	rangeCal := page.Locator("#range-cal")
	day12 := rangeCal.Locator("button", pw.LocatorLocatorOptions{
		HasText: "12",
	}).First()
	if err := day12.Click(); err != nil {
		t.Fatalf("click day 12: %v", err)
	}

	// Badge should now show a date, not "none".
	if err := badge.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("badge wait: %v", err)
	}

	newText, err := badge.TextContent()
	if err != nil {
		t.Fatalf("badge text after click: %v", err)
	}
	if newText == "none" || newText == "" {
		t.Errorf("badge did not update after clicking day 12, got %q", newText)
	}
}
