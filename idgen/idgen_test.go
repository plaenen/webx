package idgen

import (
	"strings"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	prefix := "test"
	id := Generate(prefix)

	if !strings.HasPrefix(id, prefix+"_") {
		t.Errorf("expected id to start with prefix %q, got %q", prefix+"_", id)
	}

	parts := strings.Split(id, "_")
	if len(parts) != 2 {
		t.Fatalf("expected id to have 2 parts separated by '_', got %d", len(parts))
	}

	ulidPart := parts[1]
	if len(ulidPart) != 26 {
		t.Errorf("expected ULID part to be 26 characters long, got %d", len(ulidPart))
	}

	if ulidPart != strings.ToLower(ulidPart) {
		t.Errorf("expected ULID part to be lowercase, got %q", ulidPart)
	}
}

func TestNewConcurrency(t *testing.T) {
	const numGoroutines = 100
	const numIDsPerGoroutine = 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	ids := make(chan string, numGoroutines*numIDsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIDsPerGoroutine; j++ {
				ids <- Generate("test")
			}
		}()
	}

	wg.Wait()
	close(ids)

	uniqueIDs := make(map[string]struct{})
	for id := range ids {
		if _, exists := uniqueIDs[id]; exists {
			t.Errorf("duplicate ID found: %s", id)
		}
		uniqueIDs[id] = struct{}{}
	}
}

func TestMonotonicity(t *testing.T) {
	// Best-effort check for monotonicity.
	// Since we are running in a test environment, we might not hit the same millisecond often enough to trigger the monotonic increment
	// unless we generate IDs very fast.
	const count = 10000
	ids := make([]string, count)

	for i := 0; i < count; i++ {
		ids[i] = Generate("mono")
	}

	for i := 1; i < count; i++ {
		prev := ids[i-1]
		curr := ids[i]
		if curr <= prev {
			// Extract ULID part for more detailed error message if needed, but simple string comparison works for time-ordered IDs
			t.Errorf("IDs are not strictly increasing: %s followed by %s", prev, curr)
		}
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Generate("bench")
	}
}
