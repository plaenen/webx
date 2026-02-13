package moneyinput

import (
	"math"
	"testing"
)

func TestParseAmount(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr string
	}{
		{"empty", "", 0, ""},
		{"whitespace", "  ", 0, ""},
		{"integer", "100", 100, ""},
		{"decimal", "123.45", 123.45, ""},
		{"thousands", "1,234.56", 1234.56, ""},
		{"5k", "5k", 5000, ""},
		{"5K", "5K", 5000, ""},
		{"1.5M", "1.5M", 1500000, ""},
		{"1.5m", "1.5m", 1500000, ""},
		{"2B", "2B", 2000000000, ""},
		{"2.5b", "2.5b", 2500000000, ""},
		{"negative", "-500", -500, ""},
		{"negative_k", "-5k", -5000, ""},
		{"bad_input", "abc", 0, "Invalid number"},
		{"just_suffix", "k", 0, "Invalid number"},
		{"multiple_dots", "1.2.3", 0, "Invalid number"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseAmount(tt.input)
			if tt.wantErr != "" {
				if got.Error != tt.wantErr {
					t.Errorf("ParseAmount(%q).Error = %q, want %q", tt.input, got.Error, tt.wantErr)
				}
				if got.Valid {
					t.Errorf("ParseAmount(%q).Valid = true, want false", tt.input)
				}
				return
			}
			if !got.Valid {
				t.Errorf("ParseAmount(%q).Valid = false, Error = %q", tt.input, got.Error)
				return
			}
			if math.Abs(got.Value-tt.want) > 0.001 {
				t.Errorf("ParseAmount(%q).Value = %f, want %f", tt.input, got.Value, tt.want)
			}
		})
	}
}

func TestParseMoney(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		allowed []string
		wantVal float64
		wantCur string
		wantErr string
	}{
		{"empty", "", nil, 0, "", ""},
		{"amount_only", "100", nil, 100, "", ""},
		{"prefix_currency", "USD 100", nil, 100, "USD", ""},
		{"suffix_currency", "100 EUR", nil, 100, "EUR", ""},
		{"prefix_with_k", "GBP 5k", nil, 5000, "GBP", ""},
		{"suffix_with_M", "1.5M JPY", nil, 1500000, "JPY", ""},
		{"lowercase_currency", "usd 100", nil, 100, "USD", ""},
		{"allowed_ok", "USD 100", []string{"USD", "EUR"}, 100, "USD", ""},
		{"allowed_fail", "GBP 100", []string{"USD", "EUR"}, 0, "", "Currency GBP is not allowed"},
		{"too_many_parts", "USD 100 extra", nil, 0, "", "Invalid format: too many parts"},
		{"bad_amount", "USD abc", nil, 0, "", "Invalid number"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseMoney(tt.input, tt.allowed)
			if tt.wantErr != "" {
				if got.Error != tt.wantErr {
					t.Errorf("ParseMoney(%q).Error = %q, want %q", tt.input, got.Error, tt.wantErr)
				}
				return
			}
			if !got.Valid {
				t.Errorf("ParseMoney(%q).Valid = false, Error = %q", tt.input, got.Error)
				return
			}
			if math.Abs(got.Value-tt.wantVal) > 0.001 {
				t.Errorf("ParseMoney(%q).Value = %f, want %f", tt.input, got.Value, tt.wantVal)
			}
			if got.Currency != tt.wantCur {
				t.Errorf("ParseMoney(%q).Currency = %q, want %q", tt.input, got.Currency, tt.wantCur)
			}
		})
	}
}

func TestFormatAmount(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{0, "0.00"},
		{1234.56, "1,234.56"},
		{1000000, "1,000,000.00"},
		{5000, "5,000.00"},
		{99.9, "99.90"},
		{1500000, "1,500,000.00"},
		{-5000, "-5,000.00"},
		{0.5, "0.50"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := FormatAmount(tt.input)
			if got != tt.want {
				t.Errorf("FormatAmount(%f) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
