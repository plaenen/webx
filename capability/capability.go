package capability

import "strings"

// Set is a collection of capability strings in "resource:action" format.
// Supports exact match, resource wildcard ("resource:*"), and superadmin ("*").
type Set []string

// Can reports whether this Set satisfies the required capability.
func (s Set) Can(required string) bool {
	for _, c := range s {
		if c == "*" || c == required {
			return true
		}
		if prefix, ok := strings.CutSuffix(c, ":*"); ok {
			if res, _, ok := strings.Cut(required, ":"); ok && res == prefix {
				return true
			}
		}
	}
	return false
}

// CanAny reports whether this Set satisfies at least one of the required capabilities.
func (s Set) CanAny(required ...string) bool {
	for _, r := range required {
		if s.Can(r) {
			return true
		}
	}
	return false
}

// CanAll reports whether this Set satisfies all of the required capabilities.
func (s Set) CanAll(required ...string) bool {
	for _, r := range required {
		if !s.Can(r) {
			return false
		}
	}
	return true
}
