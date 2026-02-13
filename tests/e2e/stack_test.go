package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestStackPage_StacksRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/stack", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	stacks := page.Locator(".stack")
	count, err := stacks.Count()
	if err != nil {
		t.Fatalf("count stacks: %v", err)
	}
	if count == 0 {
		t.Error("no stack components found on stack page")
	}
}

func TestStackPage_BasicStackHasChildren(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/stack", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	basic := page.Locator("#stack-basic")
	children := basic.Locator("> div")
	count, err := children.Count()
	if err != nil {
		t.Fatalf("count children: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 children in basic stack, got %d", count)
	}
}

func TestStackPage_PositionVariantsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/stack", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	positions := []struct {
		id    string
		class string
	}{
		{"#stack-top", "stack-top"},
		{"#stack-bottom", "stack-bottom"},
		{"#stack-start", "stack-start"},
		{"#stack-end", "stack-end"},
	}
	for _, p := range positions {
		loc := page.Locator(p.id)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", p.id, err)
		}
		if count == 0 {
			t.Errorf("no element found with id %s", p.id)
			continue
		}
		cls, err := loc.GetAttribute("class")
		if err != nil {
			t.Fatalf("get class for %s: %v", p.id, err)
		}
		found := false
		for _, c := range splitClasses(cls) {
			if c == p.class {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected %s to have class %q, got %q", p.id, p.class, cls)
		}
	}
}

func splitClasses(s string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ' ' {
			if i > start {
				result = append(result, s[start:i])
			}
			start = i + 1
		}
	}
	return result
}
