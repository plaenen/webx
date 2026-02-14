package memory

import (
	"fmt"
	"sync"

	"github.com/plaenen/webx"
)

var _ webx.SubscriptionStore = (*Store)(nil)

// Store is an in-memory SubscriptionStore.
type Store struct {
	mu   sync.RWMutex
	data map[string]webx.Subscription
}

// New creates an in-memory SubscriptionStore.
func New() *Store {
	return &Store{data: make(map[string]webx.Subscription)}
}

func (s *Store) Create(sub webx.Subscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[sub.ID]; ok {
		return fmt.Errorf("subscription %q already exists", sub.ID)
	}
	s.data[sub.ID] = sub
	return nil
}

func (s *Store) Get(id string) (*webx.Subscription, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sub, ok := s.data[id]
	if !ok {
		return nil, nil
	}
	return &sub, nil
}

func (s *Store) GetByWorkspace(workspaceID string) (*webx.Subscription, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, sub := range s.data {
		if sub.WorkspaceID == workspaceID {
			return &sub, nil
		}
	}
	return nil, nil
}

func (s *Store) Update(sub webx.Subscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[sub.ID]; !ok {
		return fmt.Errorf("subscription %q not found", sub.ID)
	}
	s.data[sub.ID] = sub
	return nil
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, id)
	return nil
}
