package money

import (
	"strings"

	"github.com/shopspring/decimal"
)

// P returns a pointer to an int.
// Useful for setting the Precision prop in Money component.
// Example: Money(Props{Precision: money.P(0)})
func P(v int) *int {
	return &v
}

// FormatDecimal formats a decimal with thousands separators and fixed precision.
// Example: 1234.56, 2 -> "1,234.56"
func FormatDecimal(d decimal.Decimal, precision int) string {
	fixed := d.StringFixed(int32(precision))
	parts := strings.Split(fixed, ".")
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = parts[1]
	}

	// Insert commas in intPart
	var result []byte
	// Handle negative sign
	start := 0
	if len(intPart) > 0 && intPart[0] == '-' {
		result = append(result, '-')
		start = 1
	}

	rawInt := intPart[start:]
	n := len(rawInt)
	for i, c := range rawInt {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}

	if precision > 0 {
		return string(result) + "." + decPart
	}
	return string(result)
}
