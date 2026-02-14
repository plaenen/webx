package webx

import (
	"context"

	"github.com/plaenen/webx/capability"
)

// WorkspaceContext holds the resolved workspace, membership, role, and
// capabilities for the current request.
type WorkspaceContext struct {
	Workspace    *Workspace
	Membership   *Membership
	Role         *Role
	Capabilities capability.Set
}

type wsCtxKey struct{}

// WithWorkspaceContext returns a new context with the WorkspaceContext attached.
func WithWorkspaceContext(ctx context.Context, wsc *WorkspaceContext) context.Context {
	return context.WithValue(ctx, wsCtxKey{}, wsc)
}

// GetWorkspaceContext returns the WorkspaceContext from the context, or nil.
func GetWorkspaceContext(ctx context.Context) *WorkspaceContext {
	wsc, _ := ctx.Value(wsCtxKey{}).(*WorkspaceContext)
	return wsc
}
