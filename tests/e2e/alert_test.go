package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestAlertPage_AlertsRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/alert", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	alerts := page.Locator("[role='alert']")
	count, err := alerts.Count()
	if err != nil {
		t.Fatalf("count alerts: %v", err)
	}
	if count == 0 {
		t.Error("no alert components found on alert page")
	}
}

func TestAlertPage_VariantsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/alert", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	variants := []string{
		"alert-info",
		"alert-success",
		"alert-warning",
		"alert-error",
	}
	for _, v := range variants {
		loc := page.Locator("." + v)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", v, err)
		}
		if count == 0 {
			t.Errorf("no %s alert found", v)
		}
	}
}

func TestAlertPage_StylesPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/alert", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	styles := []string{
		"alert-outline",
		"alert-dash",
		"alert-soft",
	}
	for _, s := range styles {
		loc := page.Locator("." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s alert found", s)
		}
	}
}

func TestAlertPage_HasRoleAlert(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/alert", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Every .alert element should have role="alert" for accessibility.
	alertsWithRole := page.Locator(".alert[role='alert']")
	allAlerts := page.Locator(".alert")

	withRole, err := alertsWithRole.Count()
	if err != nil {
		t.Fatalf("count alerts with role: %v", err)
	}
	total, err := allAlerts.Count()
	if err != nil {
		t.Fatalf("count all alerts: %v", err)
	}
	if withRole != total {
		t.Errorf("%d of %d alerts have role='alert'", withRole, total)
	}
}
