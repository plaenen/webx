package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestKbdPage_KbdsRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/kbd", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	kbds := page.Locator("main kbd.kbd")
	count, err := kbds.Count()
	if err != nil {
		t.Fatalf("count kbds: %v", err)
	}
	if count == 0 {
		t.Error("no kbd components found on kbd page")
	}
}

func TestKbdPage_SizeVariants(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/kbd", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	sizes := []string{"kbd-xs", "kbd-sm", "kbd-md", "kbd-lg", "kbd-xl"}
	for _, s := range sizes {
		loc := page.Locator("main kbd." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s kbd found", s)
		}
	}
}

func TestKbdPage_SemanticElement(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/kbd", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// All kbd components should use the semantic <kbd> element
	kbdElements := page.Locator("main kbd")
	nonKbdClass := page.Locator("main .kbd:not(kbd)")

	kbdCount, err := kbdElements.Count()
	if err != nil {
		t.Fatalf("count kbd elements: %v", err)
	}
	nonKbdCount, err := nonKbdClass.Count()
	if err != nil {
		t.Fatalf("count non-kbd elements: %v", err)
	}

	if kbdCount == 0 {
		t.Error("no <kbd> elements found")
	}
	if nonKbdCount > 0 {
		t.Errorf("found %d elements with .kbd class that are not <kbd> elements", nonKbdCount)
	}
}
