package session

import "context"

type contextKey struct{}

// WithUserSession returns a new context with the UserSession attached.
func WithUserSession(ctx context.Context, s *UserSession) context.Context {
	return context.WithValue(ctx, contextKey{}, s)
}

// GetUserSession returns the UserSession from the context, or nil if not present.
func GetUserSession(ctx context.Context) *UserSession {
	s, _ := ctx.Value(contextKey{}).(*UserSession)
	return s
}
