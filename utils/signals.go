package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SignalManager provides a structured way to manage Datastar signals.
// It namespaces signals by ID so multiple instances of the same component
// on a page don't collide.
type SignalManager struct {
	ID          string
	Signals     any
	DataSignals string
}

// Signals creates a new SignalManager with the given ID and initial state.
// The ID is sanitized (hyphens → underscores) for JavaScript compatibility.
// The signalsStruct should have json tags for each property.
func Signals(id string, signalsStruct any) *SignalManager {
	sanitizedID := strings.ReplaceAll(id, "-", "_")

	nested := map[string]any{
		sanitizedID: signalsStruct,
	}

	jsonBytes, err := json.Marshal(nested)
	if err != nil {
		jsonBytes = []byte("{}")
	}

	return &SignalManager{
		ID:          sanitizedID,
		Signals:     signalsStruct,
		DataSignals: string(jsonBytes),
	}
}

// Signal returns a reference to a signal property: "$componentID.property"
func (sm *SignalManager) Signal(property string) string {
	return fmt.Sprintf("$%s.%s", sm.ID, property)
}

// Toggle returns an expression that flips a boolean signal.
func (sm *SignalManager) Toggle(property string) string {
	ref := sm.Signal(property)
	return fmt.Sprintf("%s = !%s", ref, ref)
}

// Set returns an assignment expression for a signal property.
func (sm *SignalManager) Set(property, value string) string {
	return fmt.Sprintf("%s = %s", sm.Signal(property), value)
}

// SetString returns an assignment expression with proper JS string quoting.
func (sm *SignalManager) SetString(property, value string) string {
	return fmt.Sprintf("%s = '%s'", sm.Signal(property), value)
}

// Equals returns a strict equality comparison expression.
func (sm *SignalManager) Equals(property, value string) string {
	return fmt.Sprintf("%s === '%s'", sm.Signal(property), value)
}

// NotEquals returns a strict inequality comparison expression.
func (sm *SignalManager) NotEquals(property, value string) string {
	return fmt.Sprintf("%s !== '%s'", sm.Signal(property), value)
}

// Conditional returns a ternary expression based on a signal property.
func (sm *SignalManager) Conditional(property, trueValue, falseValue string) string {
	return fmt.Sprintf("%s ? %s : %s", sm.Signal(property), trueValue, falseValue)
}

// ConditionalAction executes an action only when a condition is true.
func (sm *SignalManager) ConditionalAction(condition, property, value string) string {
	return fmt.Sprintf("%s ? (%s) : void 0", condition, sm.Set(property, value))
}
