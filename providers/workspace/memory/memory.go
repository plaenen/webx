package memory

import (
	"fmt"
	"sync"

	"github.com/plaenen/webx"
)

var _ webx.WorkspaceStore = (*Store)(nil)

// Store is an in-memory WorkspaceStore.
type Store struct {
	mu   sync.RWMutex
	data map[string]webx.Workspace
}

// New creates an in-memory WorkspaceStore.
func New() *Store {
	return &Store{data: make(map[string]webx.Workspace)}
}

func (s *Store) Create(ws webx.Workspace) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[ws.ID]; ok {
		return fmt.Errorf("workspace %q already exists", ws.ID)
	}
	s.data[ws.ID] = ws
	return nil
}

func (s *Store) Get(id string) (*webx.Workspace, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ws, ok := s.data[id]
	if !ok {
		return nil, nil
	}
	return &ws, nil
}

func (s *Store) List() ([]webx.Workspace, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]webx.Workspace, 0, len(s.data))
	for _, ws := range s.data {
		result = append(result, ws)
	}
	return result, nil
}

func (s *Store) Update(ws webx.Workspace) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[ws.ID]; !ok {
		return fmt.Errorf("workspace %q not found", ws.ID)
	}
	s.data[ws.ID] = ws
	return nil
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, id)
	return nil
}
