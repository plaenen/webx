package memory

import (
	"fmt"
	"sync"

	"github.com/plaenen/webx"
)

// Store is an in-memory SessionStore backed by a sync.RWMutex-protected map.
// Suitable for development and single-instance deployments.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

// New returns a ready-to-use in-memory session store.
func New() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

// Compile-time check that Store implements webx.SessionStore.
var _ webx.SessionStore = (*Store)(nil)

func sessionKey(sessionID, key string) string {
	return fmt.Sprintf("%s:%s", sessionID, key)
}

func (s *Store) Get(sessionID string, key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[sessionKey(sessionID, key)]
	if !ok {
		return "", nil
	}
	return v, nil
}

func (s *Store) Set(sessionID string, key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[sessionKey(sessionID, key)] = value
	return nil
}

func (s *Store) Delete(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	prefix := sessionID + ":"
	for k := range s.data {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			delete(s.data, k)
		}
	}
	return nil
}
