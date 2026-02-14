package memory

import (
	"fmt"
	"sync"

	"github.com/plaenen/webx"
)

var _ webx.MembershipStore = (*Store)(nil)

// Store is an in-memory MembershipStore.
type Store struct {
	mu   sync.RWMutex
	data map[string]webx.Membership
}

// New creates an in-memory MembershipStore.
func New() *Store {
	return &Store{data: make(map[string]webx.Membership)}
}

func (s *Store) Create(m webx.Membership) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[m.ID]; ok {
		return fmt.Errorf("membership %q already exists", m.ID)
	}
	s.data[m.ID] = m
	return nil
}

func (s *Store) Get(id string) (*webx.Membership, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	m, ok := s.data[id]
	if !ok {
		return nil, nil
	}
	return &m, nil
}

func (s *Store) GetByUserAndWorkspace(userID, workspaceID string) (*webx.Membership, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, m := range s.data {
		if m.UserID == userID && m.WorkspaceID == workspaceID {
			return &m, nil
		}
	}
	return nil, nil
}

func (s *Store) ListByWorkspace(workspaceID string) ([]webx.Membership, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []webx.Membership
	for _, m := range s.data {
		if m.WorkspaceID == workspaceID {
			result = append(result, m)
		}
	}
	return result, nil
}

func (s *Store) ListByUser(userID string) ([]webx.Membership, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []webx.Membership
	for _, m := range s.data {
		if m.UserID == userID {
			result = append(result, m)
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
