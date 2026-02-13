package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestFileUploadPage_InputRendersFileInput(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/file-upload", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	inputs := page.Locator("input[type='file']")
	count, err := inputs.Count()
	if err != nil {
		t.Fatalf("count file inputs: %v", err)
	}
	if count == 0 {
		t.Error("no file input elements found on file upload page")
	}
}

func TestFileUploadPage_MultipleAttributePresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/file-upload", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Use attribute selector to verify multiple is present
	multiInput := page.Locator("#upload-multi-input[multiple]")
	count, err := multiInput.Count()
	if err != nil {
		t.Fatalf("count multi input: %v", err)
	}
	if count == 0 {
		t.Error("expected multiple attribute on multi file input")
	}
}

func TestFileUploadPage_AcceptAttributePresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/file-upload", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	restrictedInput := page.Locator("#upload-restricted-input")
	count, err := restrictedInput.Count()
	if err != nil {
		t.Fatalf("count restricted input: %v", err)
	}
	if count == 0 {
		t.Fatal("restricted file input not found")
	}
	accept, err := restrictedInput.GetAttribute("accept")
	if err != nil {
		t.Fatalf("get accept attr: %v", err)
	}
	if accept != "image/*" {
		t.Errorf("expected accept='image/*', got %q", accept)
	}
}

func TestFileUploadPage_FileListContainerExists(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/file-upload", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	listDiv := page.Locator("#upload-multi-list")
	count, err := listDiv.Count()
	if err != nil {
		t.Fatalf("count list div: %v", err)
	}
	if count == 0 {
		t.Error("file list container #upload-multi-list not found")
	}
}
