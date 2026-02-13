package moneyinput

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// ParsedAmount holds the result of parsing a numeric string.
type ParsedAmount struct {
	Value float64
	Valid bool
	Error string
}

// ParsedMoney holds the result of parsing a money string.
type ParsedMoney struct {
	Value    float64
	Currency string
	Valid    bool
	Error    string
}

var shorthandMultipliers = map[byte]float64{
	'k': 1_000,
	'K': 1_000,
	'm': 1_000_000,
	'M': 1_000_000,
	'b': 1_000_000_000,
	'B': 1_000_000_000,
}

var currencyCodeRegex = regexp.MustCompile(`^[A-Z]{3}$`)

// ParseAmount parses a string into a float64 amount.
// Handles plain numbers, thousands separators, and shorthand suffixes (k, M, B).
// Returns Valid=true with zero Value for empty input.
func ParseAmount(raw string) ParsedAmount {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ParsedAmount{Valid: true}
	}

	var multiplier float64 = 1
	if len(s) > 0 {
		last := s[len(s)-1]
		if m, ok := shorthandMultipliers[last]; ok {
			multiplier = m
			s = s[:len(s)-1]
		}
	}

	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)

	if s == "" || s == "-" || s == "+" {
		return ParsedAmount{Error: "Invalid number"}
	}

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return ParsedAmount{Error: "Invalid number"}
	}

	if math.IsInf(val, 0) || math.IsNaN(val) {
		return ParsedAmount{Error: "Invalid number"}
	}

	val *= multiplier

	return ParsedAmount{Value: val, Valid: true}
}

// ParseMoney parses a string that may contain a currency code and amount.
// Accepts "USD 100", "100 EUR", or plain "100".
// If allowedCurrencies is non-empty, the detected currency must be in the list.
func ParseMoney(raw string, allowedCurrencies []string) ParsedMoney {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ParsedMoney{Valid: true}
	}

	parts := strings.Fields(s)
	var amountStr, currency string

	switch len(parts) {
	case 1:
		amountStr = parts[0]
	case 2:
		if currencyCodeRegex.MatchString(strings.ToUpper(parts[0])) {
			currency = strings.ToUpper(parts[0])
			amountStr = parts[1]
		} else if currencyCodeRegex.MatchString(strings.ToUpper(parts[1])) {
			currency = strings.ToUpper(parts[1])
			amountStr = parts[0]
		} else {
			return ParsedMoney{Error: "Invalid format: expected amount or currency + amount"}
		}
	default:
		return ParsedMoney{Error: "Invalid format: too many parts"}
	}

	parsed := ParseAmount(amountStr)
	if !parsed.Valid {
		return ParsedMoney{Error: parsed.Error}
	}

	if currency != "" && len(allowedCurrencies) > 0 {
		allowed := false
		for _, c := range allowedCurrencies {
			if strings.EqualFold(c, currency) {
				allowed = true
				break
			}
		}
		if !allowed {
			return ParsedMoney{Error: fmt.Sprintf("Currency %s is not allowed", currency)}
		}
	}

	return ParsedMoney{
		Value:    parsed.Value,
		Currency: currency,
		Valid:    true,
	}
}

// FormatAmount formats a float64 with 2 decimal places and thousands separators.
func FormatAmount(val float64) string {
	s := fmt.Sprintf("%.2f", val)

	dotIdx := strings.Index(s, ".")
	intPart := s[:dotIdx]
	decPart := s[dotIdx:]

	negative := false
	if intPart[0] == '-' {
		negative = true
		intPart = intPart[1:]
	}

	if len(intPart) > 3 {
		var b strings.Builder
		remainder := len(intPart) % 3
		if remainder > 0 {
			b.WriteString(intPart[:remainder])
		}
		for i := remainder; i < len(intPart); i += 3 {
			if b.Len() > 0 {
				b.WriteByte(',')
			}
			b.WriteString(intPart[i : i+3])
		}
		intPart = b.String()
	}

	if negative {
		return "-" + intPart + decPart
	}
	return intPart + decPart
}
