package memory

import (
	"fmt"
	"sync"
	"time"

	"github.com/plaenen/webx"
	"github.com/plaenen/webx/idgen"
)

type MemTokenStore struct {
	mu     sync.RWMutex
	tokens map[string]*webx.Token
	ttl    time.Duration
}

func NewMemTokenStore(ttl time.Duration) *MemTokenStore {
	s := &MemTokenStore{
		tokens: make(map[string]*webx.Token),
		ttl:    ttl,
	}
	go s.cleanup()
	return s
}

func (s *MemTokenStore) Create(email string) (*webx.Token, error) {
	tokenVal, err := idgen.Token(32)
	if err != nil {
		return nil, fmt.Errorf("generating token: %w", err)
	}

	token := &webx.Token{
		Value:     tokenVal,
		Email:     email,
		Verified:  false,
		ExpiresAt: time.Now().Add(s.ttl),
		CreatedAt: time.Now(),
	}

	s.mu.Lock()
	s.tokens[token.Value] = token
	s.mu.Unlock()

	return token, nil
}

func (s *MemTokenStore) Get(value string) (*webx.Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	token, ok := s.tokens[value]
	if !ok {
		return nil, fmt.Errorf("token not found")
	}
	if token.Expired() {
		return nil, fmt.Errorf("token expired")
	}
	return token, nil
}

func (s *MemTokenStore) MarkVerified(value string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if t, ok := s.tokens[value]; ok && !t.Expired() {
		t.Verified = true
		return true
	}
	return false
}

func (s *MemTokenStore) Delete(value string) {
	s.mu.Lock()
	delete(s.tokens, value)
	s.mu.Unlock()
}

func (s *MemTokenStore) ListAll() ([]*webx.Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*webx.Token
	for _, t := range s.tokens {
		if !t.Expired() {
			result = append(result, t)
		}
	}
	return result, nil
}

func (s *MemTokenStore) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		s.mu.Lock()
		for k, t := range s.tokens {
			if t.Expired() {
				delete(s.tokens, k)
			}
		}
		s.mu.Unlock()
	}
}
