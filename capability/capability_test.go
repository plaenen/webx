package capability_test

import (
	"testing"

	"github.com/plaenen/webx/capability"
)

func TestCan_ExactMatch(t *testing.T) {
	s := capability.Set{"invoices:read", "users:write"}

	if !s.Can("invoices:read") {
		t.Error("expected invoices:read to match")
	}
	if !s.Can("users:write") {
		t.Error("expected users:write to match")
	}
	if s.Can("invoices:write") {
		t.Error("expected invoices:write not to match")
	}
}

func TestCan_ResourceWildcard(t *testing.T) {
	s := capability.Set{"invoices:*"}

	if !s.Can("invoices:read") {
		t.Error("expected invoices:* to match invoices:read")
	}
	if !s.Can("invoices:write") {
		t.Error("expected invoices:* to match invoices:write")
	}
	if !s.Can("invoices:delete") {
		t.Error("expected invoices:* to match invoices:delete")
	}
	if s.Can("users:read") {
		t.Error("expected invoices:* not to match users:read")
	}
}

func TestCan_Superadmin(t *testing.T) {
	s := capability.Set{"*"}

	if !s.Can("invoices:read") {
		t.Error("expected * to match invoices:read")
	}
	if !s.Can("users:write") {
		t.Error("expected * to match users:write")
	}
	if !s.Can("anything:at:all") {
		t.Error("expected * to match anything")
	}
}

func TestCan_EmptySet(t *testing.T) {
	s := capability.Set{}

	if s.Can("invoices:read") {
		t.Error("expected empty set to match nothing")
	}
}

func TestCan_NilSet(t *testing.T) {
	var s capability.Set

	if s.Can("invoices:read") {
		t.Error("expected nil set to match nothing")
	}
}

func TestCan_EmptyRequired(t *testing.T) {
	s := capability.Set{"invoices:read"}

	if s.Can("") {
		t.Error("expected empty required not to match")
	}
}

func TestCan_WildcardDoesNotMatchPlainString(t *testing.T) {
	s := capability.Set{"invoices:*"}

	// "invoices:*" should not match a capability without a colon
	if s.Can("invoices") {
		t.Error("expected invoices:* not to match bare 'invoices'")
	}
}

func TestCanAny(t *testing.T) {
	s := capability.Set{"invoices:read", "reports:read"}

	if !s.CanAny("invoices:read", "users:write") {
		t.Error("expected CanAny to match when at least one matches")
	}
	if s.CanAny("users:write", "billing:read") {
		t.Error("expected CanAny not to match when none match")
	}
}

func TestCanAll(t *testing.T) {
	s := capability.Set{"invoices:read", "reports:read", "users:write"}

	if !s.CanAll("invoices:read", "reports:read") {
		t.Error("expected CanAll to match when all match")
	}
	if s.CanAll("invoices:read", "billing:read") {
		t.Error("expected CanAll not to match when one is missing")
	}
}

func TestCanAll_Empty(t *testing.T) {
	s := capability.Set{"invoices:read"}

	// CanAll with no args should return true (vacuous truth).
	if !s.CanAll() {
		t.Error("expected CanAll with no args to return true")
	}
}

func TestCanAny_Empty(t *testing.T) {
	s := capability.Set{"invoices:read"}

	// CanAny with no args should return false (nothing to match).
	if s.CanAny() {
		t.Error("expected CanAny with no args to return false")
	}
}
