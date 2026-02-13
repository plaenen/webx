package e2e_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

func TestChatPage_ChatsRender(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/chat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	chats := page.Locator(".chat")
	count, err := chats.Count()
	if err != nil {
		t.Fatalf("count chats: %v", err)
	}
	if count == 0 {
		t.Error("no chat components found on chat page")
	}
}

func TestChatPage_PositionsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/chat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	positions := []string{"chat-start", "chat-end"}
	for _, p := range positions {
		loc := page.Locator("." + p)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", p, err)
		}
		if count == 0 {
			t.Errorf("no %s chat found", p)
		}
	}
}

func TestChatPage_BubbleColorsPresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/chat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	colors := []string{
		"chat-bubble-primary",
		"chat-bubble-secondary",
		"chat-bubble-accent",
		"chat-bubble-neutral",
		"chat-bubble-info",
		"chat-bubble-success",
		"chat-bubble-warning",
		"chat-bubble-error",
	}
	for _, c := range colors {
		loc := page.Locator("." + c)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", c, err)
		}
		if count == 0 {
			t.Errorf("no %s bubble found", c)
		}
	}
}

func TestChatPage_HeaderFooterImagePresent(t *testing.T) {
	page := newPage(t)
	if _, err := page.Goto(baseURL+"/components/chat", pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Fatalf("goto: %v", err)
	}

	parts := []string{"chat-header", "chat-footer", "chat-image"}
	for _, p := range parts {
		loc := page.Locator("." + p)
		count, err := loc.Count()
		if err != nil {
			t.Fatalf("count %s: %v", p, err)
		}
		if count == 0 {
			t.Errorf("no %s element found", p)
		}
	}
}
