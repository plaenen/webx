package utils

import (
	"fmt"
	"strings"
)

// DataClass builds data-class attribute values: JS objects mapping
// CSS class names to signal conditions.
type DataClass struct {
	classes []dataClassEntry
}

type dataClassEntry struct {
	className string
	condition string
}

// NewDataClass creates a new DataClass builder.
func NewDataClass() *DataClass {
	return &DataClass{}
}

// Add maps a CSS class name to a signal condition.
func (d *DataClass) Add(className, condition string) *DataClass {
	d.classes = append(d.classes, dataClassEntry{className, condition})
	return d
}

// Build returns the JS object string for the data-class attribute.
func (d *DataClass) Build() string {
	if len(d.classes) == 0 {
		return "{}"
	}

	parts := make([]string, len(d.classes))
	for i, entry := range d.classes {
		parts[i] = fmt.Sprintf("'%s': %s", entry.className, entry.condition)
	}

	return "{" + strings.Join(parts, ", ") + "}"
}
