package memory

import (
	"testing"
	"time"

	"github.com/plaenen/webx"
)

func TestCreateAndGet(t *testing.T) {
	s := New()
	p := webx.Provider{
		ID:        "prov_1",
		Name:      "AWS",
		Type:      "managed",
		CreatedAt: time.Now(),
	}

	if err := s.Create(p); err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := s.Get("prov_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got == nil || got.ID != "prov_1" {
		t.Fatalf("expected prov_1, got %v", got)
	}
}

func TestCreateDuplicate(t *testing.T) {
	s := New()
	p := webx.Provider{ID: "prov_1", CreatedAt: time.Now()}

	if err := s.Create(p); err != nil {
		t.Fatalf("first Create: %v", err)
	}
	if err := s.Create(p); err == nil {
		t.Fatal("expected error on duplicate Create")
	}
}

func TestList(t *testing.T) {
	s := New()
	_ = s.Create(webx.Provider{ID: "prov_1", CreatedAt: time.Now()})
	_ = s.Create(webx.Provider{ID: "prov_2", CreatedAt: time.Now()})

	list, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 providers, got %d", len(list))
	}
}

func TestListByType(t *testing.T) {
	s := New()
	_ = s.Create(webx.Provider{ID: "prov_1", Type: "managed", CreatedAt: time.Now()})
	_ = s.Create(webx.Provider{ID: "prov_2", Type: "selfhosted", CreatedAt: time.Now()})
	_ = s.Create(webx.Provider{ID: "prov_3", Type: "managed", CreatedAt: time.Now()})

	managed, err := s.ListByType("managed")
	if err != nil {
		t.Fatalf("ListByType: %v", err)
	}
	if len(managed) != 2 {
		t.Fatalf("expected 2 managed, got %d", len(managed))
	}

	selfhosted, err := s.ListByType("selfhosted")
	if err != nil {
		t.Fatalf("ListByType: %v", err)
	}
	if len(selfhosted) != 1 {
		t.Fatalf("expected 1 selfhosted, got %d", len(selfhosted))
	}
}

func TestListByOwner(t *testing.T) {
	s := New()
	_ = s.Create(webx.Provider{ID: "prov_1", OwnerID: "user_1", CreatedAt: time.Now()})
	_ = s.Create(webx.Provider{ID: "prov_2", OwnerID: "user_2", CreatedAt: time.Now()})

	list, err := s.ListByOwner("user_1")
	if err != nil {
		t.Fatalf("ListByOwner: %v", err)
	}
	if len(list) != 1 || list[0].ID != "prov_1" {
		t.Fatalf("expected [prov_1], got %v", list)
	}
}

func TestUpdate(t *testing.T) {
	s := New()
	p := webx.Provider{ID: "prov_1", Name: "AWS", CreatedAt: time.Now()}
	_ = s.Create(p)

	p.Name = "Amazon Web Services"
	if err := s.Update(p); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, _ := s.Get("prov_1")
	if got.Name != "Amazon Web Services" {
		t.Fatalf("expected updated name, got %q", got.Name)
	}
}

func TestUpdateNotFound(t *testing.T) {
	s := New()
	if err := s.Update(webx.Provider{ID: "nope"}); err == nil {
		t.Fatal("expected error on Update for missing provider")
	}
}

func TestDelete(t *testing.T) {
	s := New()
	_ = s.Create(webx.Provider{ID: "prov_1", CreatedAt: time.Now()})

	if err := s.Delete("prov_1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	got, _ := s.Get("prov_1")
	if got != nil {
		t.Fatal("expected nil after Delete")
	}
}
