package memory

import (
	"fmt"
	"sync"

	"github.com/plaenen/webx"
)

var _ webx.RoleStore = (*Store)(nil)

// Store is an in-memory RoleStore.
type Store struct {
	mu   sync.RWMutex
	data map[string]webx.Role
}

// New creates an in-memory RoleStore.
func New() *Store {
	return &Store{data: make(map[string]webx.Role)}
}

func (s *Store) Create(r webx.Role) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[r.ID]; ok {
		return fmt.Errorf("role %q already exists", r.ID)
	}
	s.data[r.ID] = r
	return nil
}

func (s *Store) Get(id string) (*webx.Role, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	r, ok := s.data[id]
	if !ok {
		return nil, nil
	}
	return &r, nil
}

func (s *Store) GetByNameAndWorkspace(name, workspaceID string) (*webx.Role, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, r := range s.data {
		if r.Name == name && r.WorkspaceID == workspaceID {
			return &r, nil
		}
	}
	return nil, nil
}

func (s *Store) ListByWorkspace(workspaceID string) ([]webx.Role, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []webx.Role
	for _, r := range s.data {
		if r.WorkspaceID == workspaceID {
			result = append(result, r)
		}
	}
	return result, nil
}

func (s *Store) Update(r webx.Role) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[r.ID]; !ok {
		return fmt.Errorf("role %q not found", r.ID)
	}
	s.data[r.ID] = r
	return nil
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.data[id]
	if ok && r.System {
		return fmt.Errorf("cannot delete system role %q", id)
	}
	delete(s.data, id)
	return nil
}
