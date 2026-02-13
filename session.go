package webx

// SessionStore provides session-scoped key/value storage.
// Implementations can back this with memory, NATS KV, Redis, etc.
type SessionStore interface {
	Get(sessionID string, key string) (string, error)
	Set(sessionID string, key string, value string) error
	Delete(sessionID string) error
}
