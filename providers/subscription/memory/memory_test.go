package memory

import (
	"testing"
	"time"

	"github.com/plaenen/webx"
)

func TestCreateAndGet(t *testing.T) {
	s := New()
	sub := webx.Subscription{
		ID:          "sub_1",
		WorkspaceID: "ws_1",
		CustomerID:  "cus_1",
		Tier:        "pro",
		Active:      true,
		CreatedAt:   time.Now(),
	}

	if err := s.Create(sub); err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := s.Get("sub_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got == nil || got.ID != "sub_1" {
		t.Fatalf("expected sub_1, got %v", got)
	}
}

func TestCreateDuplicate(t *testing.T) {
	s := New()
	sub := webx.Subscription{ID: "sub_1", WorkspaceID: "ws_1", CreatedAt: time.Now()}

	if err := s.Create(sub); err != nil {
		t.Fatalf("first Create: %v", err)
	}
	if err := s.Create(sub); err == nil {
		t.Fatal("expected error on duplicate Create")
	}
}

func TestGetByWorkspace(t *testing.T) {
	s := New()
	_ = s.Create(webx.Subscription{ID: "sub_1", WorkspaceID: "ws_1", CreatedAt: time.Now()})
	_ = s.Create(webx.Subscription{ID: "sub_2", WorkspaceID: "ws_2", CreatedAt: time.Now()})

	got, err := s.GetByWorkspace("ws_2")
	if err != nil {
		t.Fatalf("GetByWorkspace: %v", err)
	}
	if got == nil || got.ID != "sub_2" {
		t.Fatalf("expected sub_2, got %v", got)
	}

	got, err = s.GetByWorkspace("ws_unknown")
	if err != nil {
		t.Fatalf("GetByWorkspace unknown: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil for unknown workspace, got %v", got)
	}
}

func TestUpdate(t *testing.T) {
	s := New()
	sub := webx.Subscription{ID: "sub_1", WorkspaceID: "ws_1", Tier: "free", CreatedAt: time.Now()}
	_ = s.Create(sub)

	sub.Tier = "pro"
	if err := s.Update(sub); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, _ := s.Get("sub_1")
	if got.Tier != "pro" {
		t.Fatalf("expected PlanPro, got %v", got.Tier)
	}
}

func TestUpdateNotFound(t *testing.T) {
	s := New()
	if err := s.Update(webx.Subscription{ID: "nope"}); err == nil {
		t.Fatal("expected error on Update for missing subscription")
	}
}

func TestDelete(t *testing.T) {
	s := New()
	_ = s.Create(webx.Subscription{ID: "sub_1", WorkspaceID: "ws_1", CreatedAt: time.Now()})

	if err := s.Delete("sub_1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	got, _ := s.Get("sub_1")
	if got != nil {
		t.Fatal("expected nil after Delete")
	}
}
