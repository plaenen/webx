package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestAvatarPage_AvatarsRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/avatar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	avatars := page.Locator(".avatar")
	count, err := avatars.Count()
	if err != nil {
		t.Fatalf("count avatars: %v", err)
	}
	if count == 0 {
		t.Error("no avatar components found on avatar page")
	}
}

func TestAvatarPage_StatusPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/avatar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	statuses := []string{
		"avatar-online",
		"avatar-offline",
	}
	for _, s := range statuses {
		loc := page.Locator("." + s)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", s, err)
		}
		if count == 0 {
			t.Errorf("no %s avatar found", s)
		}
	}
}

func TestAvatarPage_PlaceholderPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/avatar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	loc := page.Locator(".avatar-placeholder")
	count, err := loc.Count()
	if err != nil {
		t.Fatalf("count placeholder: %v", err)
	}
	if count == 0 {
		t.Error("no avatar-placeholder found")
	}
}

func TestAvatarPage_GroupPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/avatar", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	loc := page.Locator(".avatar-group")
	count, err := loc.Count()
	if err != nil {
		t.Fatalf("count group: %v", err)
	}
	if count == 0 {
		t.Error("no avatar-group found")
	}
}
