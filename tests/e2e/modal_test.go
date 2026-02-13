package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestModalPage_OpenAndClose(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/modal", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Click the first "Open Modal" button (basic modal).
	openBtn := page.Locator("button", pw.PageLocatorOptions{
		HasText: "Open Modal",
	}).First()
	if err := openBtn.Click(); err != nil {
		t.Fatalf("click open: %v", err)
	}

	// The modal-box should be visible.
	modalBox := page.Locator(".modal-box").First()
	if err := modalBox.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("modal-box did not become visible: %v", err)
	}

	// Click close button.
	closeBtn := page.Locator(".modal-action .btn").First()
	if err := closeBtn.Click(); err != nil {
		t.Fatalf("click close: %v", err)
	}

	// Modal should be hidden.
	if err := modalBox.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateHidden,
	}); err != nil {
		t.Fatalf("modal did not close: %v", err)
	}
}

func TestModalPage_BackdropCloses(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/modal", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Open the backdrop modal (second "Open Modal" button).
	openBtns := page.Locator("button", pw.PageLocatorOptions{
		HasText: "Open Modal",
	})
	if err := openBtns.Nth(1).Click(); err != nil {
		t.Fatalf("click open: %v", err)
	}

	// Wait for modal to be visible.
	modalBox := page.Locator(".modal-box").Nth(1)
	if err := modalBox.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("modal-box did not become visible: %v", err)
	}

	// Click the backdrop to close. The backdrop is full-screen but behind the
	// modal-box, so click near the top-left corner to avoid hitting the box.
	backdrop := page.Locator(".modal-backdrop").First()
	if err := backdrop.Click(pw.LocatorClickOptions{
		Position: &pw.Position{X: 10, Y: 10},
		Force:    pw.Bool(true),
	}); err != nil {
		t.Fatalf("click backdrop: %v", err)
	}

	// Modal should close.
	if err := modalBox.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateHidden,
	}); err != nil {
		t.Fatalf("modal did not close after backdrop click: %v", err)
	}
}

func TestModalPage_CornerCloseButton(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/modal", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Open the corner-close modal (third "Open Modal" button).
	openBtns := page.Locator("button", pw.PageLocatorOptions{
		HasText: "Open Modal",
	})
	if err := openBtns.Nth(2).Click(); err != nil {
		t.Fatalf("click open: %v", err)
	}

	// Wait for modal to appear.
	modalBox := page.Locator(".modal-box").Nth(2)
	if err := modalBox.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateVisible,
	}); err != nil {
		t.Fatalf("modal-box did not become visible: %v", err)
	}

	// Click the corner ✕ button.
	cornerBtn := modalBox.Locator(".btn-circle")
	if err := cornerBtn.Click(); err != nil {
		t.Fatalf("click corner close: %v", err)
	}

	// Modal should close.
	if err := modalBox.WaitFor(pw.LocatorWaitForOptions{
		State: pw.WaitForSelectorStateHidden,
	}); err != nil {
		t.Fatalf("modal did not close after corner button click: %v", err)
	}
}

func TestModalPage_PositionVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/modal", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	positions := []string{"modal-top", "modal-middle", "modal-bottom"}
	for _, p := range positions {
		loc := page.Locator("main .modal." + p)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", p, err)
		}
		if count == 0 {
			t.Errorf("no %s modal found", p)
		}
	}
}
