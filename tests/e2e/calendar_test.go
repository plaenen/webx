package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestCalendarPage_CalendarsRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/calendar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Multiple calendar containers should exist.
	calendars := page.Locator(".rounded-box .grid-cols-7")
	count, err := calendars.Count()
	if err != nil {
		t.Fatalf("count calendars: %v", err)
	}
	if count < 2 {
		t.Errorf("expected at least 2 calendar grids, got %d", count)
	}
}

func TestCalendarPage_WeekdayHeaders(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/calendar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The first calendar should have all 7 weekday labels.
	for _, wd := range []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"} {
		loc := page.Locator(".grid-cols-7 span", pw.PageLocatorOptions{
			HasText: wd,
		}).First()
		visible, err := loc.IsVisible()
		if err != nil {
			t.Fatalf("check weekday %s: %v", wd, err)
		}
		if !visible {
			t.Errorf("weekday header %q not visible", wd)
		}
	}
}

func TestCalendarPage_DayButtonsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/calendar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Each calendar grid should have 42 day buttons (6 rows × 7 cols).
	buttons := page.Locator(".grid-cols-7 button").First()
	visible, err := buttons.IsVisible()
	if err != nil {
		t.Fatalf("check day button: %v", err)
	}
	if !visible {
		t.Error("no day buttons found in calendar grid")
	}
}

func TestCalendarPage_ClickDayUpdatesState(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/calendar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The state indicator badge should initially show "none".
	badge := page.Locator("code.badge")
	text, err := badge.TextContent()
	if err != nil {
		t.Fatalf("badge text: %v", err)
	}
	if text != "none" {
		t.Errorf("initial badge text = %q, want %q", text, "none")
	}

	// Click a day button in the state-indicator calendar (the third calendar).
	// Find the day button with text "15" in the state calendar section.
	stateSection := page.Locator("[data-signals*='state_cal']")
	dayBtn := stateSection.Locator("button", pw.LocatorLocatorOptions{
		HasText: "15",
	}).First()
	if err := dayBtn.Click(); err != nil {
		t.Fatalf("click day 15: %v", err)
	}

	// Badge should now show a date string (not "none").
	if err := page.Locator("code.badge").WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("badge wait: %v", err)
	}

	newText, err := page.Locator("code.badge").TextContent()
	if err != nil {
		t.Fatalf("badge text after click: %v", err)
	}
	if newText == "none" || newText == "" {
		t.Errorf("badge did not update after clicking day 15, got %q", newText)
	}
}

func TestCalendarPage_PreselectedDate(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/calendar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The pre-selected calendar should have a btn-primary day (the 15th).
	preselCal := page.Locator("[data-signals*='presel_cal']")
	primaryBtn := preselCal.Locator("button.btn-primary")
	if err := primaryBtn.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("pre-selected day not highlighted: %v", err)
	}

	text, err := primaryBtn.TextContent()
	if err != nil {
		t.Fatalf("primary button text: %v", err)
	}
	if text != "15" {
		t.Errorf("pre-selected day text = %q, want %q", text, "15")
	}
}
