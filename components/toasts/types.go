package toasts

import "github.com/plaenen/webx/ui/toast"

// Variant represents the toast type/severity.
type Variant = toast.Variant

// Re-export variants for convenience.
const (
	VariantDefault = toast.VariantDefault
	VariantSuccess = toast.VariantSuccess
	VariantError   = toast.VariantError
	VariantWarning = toast.VariantWarning
	VariantInfo    = toast.VariantInfo
)

// Position represents where toasts appear on screen.
type Position = toast.Position

// Re-export positions for convenience.
const (
	PositionTopRight     = toast.PositionTopRight
	PositionTopLeft      = toast.PositionTopLeft
	PositionTopCenter    = toast.PositionTopCenter
	PositionBottomRight  = toast.PositionBottomRight
	PositionBottomLeft   = toast.PositionBottomLeft
	PositionBottomCenter = toast.PositionBottomCenter
)

// Action represents a button action on a toast.
type Action struct {
	Label   string // Button text
	URL     string // POST endpoint to call
	Variant string // Button variant (default, destructive, outline)
}

// ToastData represents a toast's data for signals.
type ToastData struct {
	ID          string   `json:"id"`
	Variant     string   `json:"variant"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Duration    int      `json:"duration"`
	Dismissible bool     `json:"dismissible"`
	RequiresAck bool     `json:"requiresAck"`
	Actions     []Action `json:"actions,omitempty"`
}

// ContainerProps configures the toast container.
type ContainerProps struct {
	ID         string
	Position   Position
	MaxVisible int
	Class      string
}

// ToastProps configures an individual toast.
type ToastProps struct {
	ID            string
	Title         string
	Description   string
	Variant       Variant
	Duration      int
	Dismissible   bool
	ShowIcon      bool
	ShowIndicator bool
	RequiresAck   bool
	Actions       []Action
	AckURL        string
}
