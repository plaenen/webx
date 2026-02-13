package webx

import "time"

type Token struct {
	Value     string
	Email     string
	Verified  bool
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (t *Token) Expired() bool {
	return time.Now().After(t.ExpiresAt)
}

type TokenStore interface {
	Create(email string) (*Token, error)
	Get(value string) (*Token, error)
	MarkVerified(value string) bool
	Delete(value string)
	ListAll() ([]*Token, error)
}
