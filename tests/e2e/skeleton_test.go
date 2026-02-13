package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestSkeletonPage_SkeletonRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/skeleton", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main .skeleton")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count skeleton elements: %v", err)
	}
	if count == 0 {
		t.Error("no skeleton elements found on skeleton page")
	}
}

func TestSkeletonPage_MultipleShapes(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/skeleton", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main .skeleton")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count skeleton elements: %v", err)
	}
	if count < 5 {
		t.Errorf("expected at least 5 skeleton elements, got %d", count)
	}
}

func TestSkeletonPage_SkeletonText(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/skeleton", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	textSkeleton := page.Locator("main span.skeleton.skeleton-text")
	count, err := textSkeleton.Count()
	if err != nil {
		t.Fatalf("count skeleton-text: %v", err)
	}
	if count == 0 {
		t.Error("no skeleton-text elements found")
	}
}
