package idgen

import (
	"crypto/rand"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	entropy   io.Reader
	entropyMu sync.Mutex
)

func init() {
	entropy = ulid.Monotonic(rand.Reader, 0)
}

// New generates a new ULID with the given prefix.
// Format: {prefix}_{ulid_lower_case}
func New(prefix string) func() string {
	return func() string {
		entropyMu.Lock()
		id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
		entropyMu.Unlock()

		// Format: prefix_ulid
		// Using lowercase for ULID part for consistency with common prefixed ID standards (like Stripe)
		return fmt.Sprintf("%s_%s", prefix, strings.ToLower(id.String()))
	}
}

func Generate(prefix string) string {
	return New(prefix)()
}

// Token generates a cryptographically secure random token of the specified length (in bytes).
// It returns the hex-encoded string representation of the token.
func Token(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generating token: %w", err)
	}
	return fmt.Sprintf("%x", b), nil
}
