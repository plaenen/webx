package memory

import (
	"fmt"
	"sync"

	"github.com/plaenen/webx"
)

var _ webx.ProviderStore = (*Store)(nil)

// Store is an in-memory ProviderStore.
type Store struct {
	mu   sync.RWMutex
	data map[string]webx.Provider
}

// New creates an in-memory ProviderStore.
func New() *Store {
	return &Store{data: make(map[string]webx.Provider)}
}

func (s *Store) Create(p webx.Provider) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[p.ID]; ok {
		return fmt.Errorf("provider %q already exists", p.ID)
	}
	s.data[p.ID] = p
	return nil
}

func (s *Store) Get(id string) (*webx.Provider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.data[id]
	if !ok {
		return nil, nil
	}
	return &p, nil
}

func (s *Store) List() ([]webx.Provider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]webx.Provider, 0, len(s.data))
	for _, p := range s.data {
		result = append(result, p)
	}
	return result, nil
}

func (s *Store) ListByType(t webx.ProviderType) ([]webx.Provider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []webx.Provider
	for _, p := range s.data {
		if p.Type == t {
			result = append(result, p)
		}
	}
	return result, nil
}

func (s *Store) ListByOwner(ownerID string) ([]webx.Provider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []webx.Provider
	for _, p := range s.data {
		if p.OwnerID == ownerID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (s *Store) Update(p webx.Provider) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[p.ID]; !ok {
		return fmt.Errorf("provider %q not found", p.ID)
	}
	s.data[p.ID] = p
	return nil
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, id)
	return nil
}
