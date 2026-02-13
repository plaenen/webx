package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestHoverGalleryPage_Render(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/hover-gallery", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main .hover-gallery")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count hover-gallery: %v", err)
	}
	if count == 0 {
		t.Error("no hover-gallery elements found")
	}
}

func TestHoverGalleryPage_MultipleGalleries(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/hover-gallery", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	containers := page.Locator("main .hover-gallery")
	count, err := containers.Count()
	if err != nil {
		t.Fatalf("count hover-gallery: %v", err)
	}
	if count < 3 {
		t.Errorf("expected at least 3 hover-gallery elements, got %d", count)
	}
}

func TestHoverGalleryPage_ImagesPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/hover-gallery", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	images := page.Locator("main .hover-gallery img")
	count, err := images.Count()
	if err != nil {
		t.Fatalf("count images: %v", err)
	}
	if count < 4 {
		t.Errorf("expected at least 4 images in hover galleries, got %d", count)
	}
}
