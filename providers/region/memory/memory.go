package memory

import (
	"fmt"
	"sync"

	"github.com/plaenen/webx"
)

var _ webx.RegionStore = (*Store)(nil)

// Store is an in-memory RegionStore.
type Store struct {
	mu   sync.RWMutex
	data map[string]webx.Region
}

// New creates an in-memory RegionStore.
func New() *Store {
	return &Store{data: make(map[string]webx.Region)}
}

func (s *Store) Create(r webx.Region) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[r.ID]; ok {
		return fmt.Errorf("region %q already exists", r.ID)
	}
	s.data[r.ID] = r
	return nil
}

func (s *Store) Get(id string) (*webx.Region, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	r, ok := s.data[id]
	if !ok {
		return nil, nil
	}
	return &r, nil
}

func (s *Store) ListByProvider(providerID string) ([]webx.Region, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []webx.Region
	for _, r := range s.data {
		if r.ProviderID == providerID {
			result = append(result, r)
		}
	}
	return result, nil
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, id)
	return nil
}
