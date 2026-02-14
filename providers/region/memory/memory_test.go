package memory

import (
	"testing"
	"time"

	"github.com/plaenen/webx"
)

func TestCreateAndGet(t *testing.T) {
	s := New()
	r := webx.Region{
		ID:         "reg_1",
		ProviderID: "prov_1",
		Name:       "us-east-1",
		Label:      "US East (Virginia)",
		CreatedAt:  time.Now(),
	}

	if err := s.Create(r); err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := s.Get("reg_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got == nil || got.ID != "reg_1" {
		t.Fatalf("expected reg_1, got %v", got)
	}
}

func TestCreateDuplicate(t *testing.T) {
	s := New()
	r := webx.Region{ID: "reg_1", ProviderID: "prov_1", CreatedAt: time.Now()}

	if err := s.Create(r); err != nil {
		t.Fatalf("first Create: %v", err)
	}
	if err := s.Create(r); err == nil {
		t.Fatal("expected error on duplicate Create")
	}
}

func TestListByProvider(t *testing.T) {
	s := New()
	_ = s.Create(webx.Region{ID: "reg_1", ProviderID: "prov_1", CreatedAt: time.Now()})
	_ = s.Create(webx.Region{ID: "reg_2", ProviderID: "prov_1", CreatedAt: time.Now()})
	_ = s.Create(webx.Region{ID: "reg_3", ProviderID: "prov_2", CreatedAt: time.Now()})

	list, err := s.ListByProvider("prov_1")
	if err != nil {
		t.Fatalf("ListByProvider: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 regions for prov_1, got %d", len(list))
	}

	list, err = s.ListByProvider("prov_2")
	if err != nil {
		t.Fatalf("ListByProvider: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 region for prov_2, got %d", len(list))
	}

	list, err = s.ListByProvider("prov_unknown")
	if err != nil {
		t.Fatalf("ListByProvider unknown: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected 0 regions for unknown, got %d", len(list))
	}
}

func TestGetNotFound(t *testing.T) {
	s := New()
	got, err := s.Get("nope")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != nil {
		t.Fatal("expected nil for missing region")
	}
}

func TestDelete(t *testing.T) {
	s := New()
	_ = s.Create(webx.Region{ID: "reg_1", ProviderID: "prov_1", CreatedAt: time.Now()})

	if err := s.Delete("reg_1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	got, _ := s.Get("reg_1")
	if got != nil {
		t.Fatal("expected nil after Delete")
	}
}
