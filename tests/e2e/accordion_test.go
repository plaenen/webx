package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestAccordionPage_RendersItems(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/accordion", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The arrow accordion should have three collapse items.
	items := page.Locator(".collapse.collapse-arrow")
	count, err := items.Count()
	if err != nil {
		t.Fatalf("count collapse items: %v", err)
	}
	// Arrow (3) + default-open (3) + joined (3) + state (3) = 12 arrow items
	// Plus accordion uses collapse-plus, so not counted here.
	if count < 3 {
		t.Errorf("expected at least 3 collapse-arrow items, got %d", count)
	}

	// The plus accordion should have items with collapse-plus.
	plusItems := page.Locator(".collapse.collapse-plus")
	plusCount, err := plusItems.Count()
	if err != nil {
		t.Fatalf("count plus items: %v", err)
	}
	if plusCount != 3 {
		t.Errorf("expected 3 collapse-plus items, got %d", plusCount)
	}
}

func TestAccordionPage_ClickOpensItem(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/accordion", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Use the state indicator accordion for reliable testing.
	// Click the "Alpha" title.
	title := page.Locator(".collapse-title", pw.PageLocatorOptions{
		HasText: "Alpha",
	})
	if err := title.Click(); err != nil {
		t.Fatalf("click Alpha title: %v", err)
	}

	// The state badge should show "alpha".
	if err := page.Locator("code.badge", pw.PageLocatorOptions{
		HasText: "alpha",
	}).WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("badge did not change to 'alpha': %v", err)
	}
}

func TestAccordionPage_ClickAnotherSwitchesItem(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/accordion", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Open Alpha first.
	alphaTitle := page.Locator(".collapse-title", pw.PageLocatorOptions{
		HasText: "Alpha",
	})
	if err := alphaTitle.Click(); err != nil {
		t.Fatalf("click Alpha: %v", err)
	}
	if err := page.Locator("code.badge", pw.PageLocatorOptions{
		HasText: "alpha",
	}).WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("badge did not show 'alpha': %v", err)
	}

	// Click Beta — should switch.
	betaTitle := page.Locator(".collapse-title", pw.PageLocatorOptions{
		HasText: "Beta",
	})
	if err := betaTitle.Click(); err != nil {
		t.Fatalf("click Beta: %v", err)
	}
	if err := page.Locator("code.badge", pw.PageLocatorOptions{
		HasText: "beta",
	}).WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("badge did not change to 'beta': %v", err)
	}
}

func TestAccordionPage_ClickOpenItemClosesIt(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/accordion", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Open Alpha.
	alphaTitle := page.Locator(".collapse-title", pw.PageLocatorOptions{
		HasText: "Alpha",
	})
	if err := alphaTitle.Click(); err != nil {
		t.Fatalf("click Alpha: %v", err)
	}
	if err := page.Locator("code.badge", pw.PageLocatorOptions{
		HasText: "alpha",
	}).WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("badge did not show 'alpha': %v", err)
	}

	// Click Alpha again — should close (toggle).
	if err := alphaTitle.Click(); err != nil {
		t.Fatalf("click Alpha again: %v", err)
	}
	if err := page.Locator("code.badge", pw.PageLocatorOptions{
		HasText: "none",
	}).WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("badge did not change to 'none': %v", err)
	}
}

func TestAccordionPage_DefaultOpenItem(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/accordion", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The default-open accordion has "Second Item (default open)" initially open.
	// Its parent collapse div should have the collapse-open class.
	openItem := page.Locator(".collapse-open", pw.PageLocatorOptions{
		Has: page.Locator(".collapse-title", pw.PageLocatorOptions{
			HasText: "Second Item (default open)",
		}),
	})
	if err := openItem.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("default open item did not have collapse-open class: %v", err)
	}
}
