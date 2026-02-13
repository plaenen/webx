package session

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
)

type MemSessionProvider struct {
	store map[string]UserSession
	mu    sync.Mutex
}

func NewMemSessionProvider() *MemSessionProvider {
	return &MemSessionProvider{
		store: make(map[string]UserSession),
	}
}

func (m *MemSessionProvider) Create(data UserSession) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id, err := randomID()
	if err != nil {
		return "", fmt.Errorf("generating session id: %w", err)
	}

	m.store[id] = data
	return id, nil
}

func (m *MemSessionProvider) Get(id string) (*UserSession, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.store[id]
	if !ok {
		return nil, nil
	}
	return &s, nil
}

func (m *MemSessionProvider) Update(id string, data UserSession) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.store[id]; !ok {
		return fmt.Errorf("session %q not found", id)
	}
	m.store[id] = data
	return nil
}

func (m *MemSessionProvider) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.store, id)
	return nil
}

func (m *MemSessionProvider) DeleteByEmail(email string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	count := 0
	for id, s := range m.store {
		if s.Email == email {
			delete(m.store, id)
			count++
		}
	}
	return count, nil
}

func (m *MemSessionProvider) ListAll() (map[string]UserSession, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make(map[string]UserSession, len(m.store))
	for id, s := range m.store {
		result[id] = s
	}
	return result, nil
}

func randomID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
