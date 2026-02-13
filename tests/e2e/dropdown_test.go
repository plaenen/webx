package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestDropdownPage_DropdownsRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/dropdown", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	dropdowns := page.Locator(".dropdown")
	count, err := dropdowns.Count()
	if err != nil {
		t.Fatalf("count dropdowns: %v", err)
	}
	if count == 0 {
		t.Error("no dropdown components found on dropdown page")
	}
}

func TestDropdownPage_PositionsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/dropdown", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	positions := []string{
		"dropdown-top",
		"dropdown-bottom",
		"dropdown-left",
		"dropdown-right",
	}
	for _, p := range positions {
		loc := page.Locator("." + p)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", p, err)
		}
		if count == 0 {
			t.Errorf("no %s dropdown found", p)
		}
	}
}

func TestDropdownPage_ClickOpensDropdown(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/dropdown", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Click the basic dropdown trigger.
	trigger := page.Locator("#dd-basic summary")
	if err := trigger.Click(); err != nil {
		t.Fatalf("click trigger: %v", err)
	}

	// The details element should have the open attribute.
	if err := page.Locator("#dd-basic[open]").WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateAttached,
	}); err != nil {
		t.Fatalf("dropdown did not open: %v", err)
	}
}

func TestDropdownPage_CloseOnItemClick(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/dropdown", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Open the "actions" dropdown.
	trigger := page.Locator("#actions summary")
	if err := trigger.Click(); err != nil {
		t.Fatalf("click actions trigger: %v", err)
	}

	if err := page.Locator("#actions[open]").WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateAttached,
	}); err != nil {
		t.Fatalf("actions dropdown did not open: %v", err)
	}

	// Click an item in the dropdown content.
	item := page.Locator("#actions .dropdown-content li a").First()
	if err := item.Click(); err != nil {
		t.Fatalf("click item: %v", err)
	}

	// The dropdown should close (open attribute removed).
	if err := page.Locator("#actions:not([open])").WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateAttached,
	}); err != nil {
		t.Fatalf("dropdown did not close after item click: %v", err)
	}
}
