package validators

import "testing"

func TestEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		checkMX bool
		valid   bool
		errMsg  string
	}{
		// Valid emails
		{name: "simple email", email: "test@example.com", valid: true},
		{name: "with subdomain", email: "test@mail.example.com", valid: true},
		{name: "with plus", email: "test+tag@example.com", valid: true},
		{name: "with dots in local", email: "first.last@example.com", valid: true},
		{name: "short TLD", email: "test@example.io", valid: true},
		{name: "long TLD", email: "test@example.museum", valid: true},
		{name: "numeric domain", email: "test@123.com", valid: true},
		{name: "hyphen in domain", email: "test@my-domain.com", valid: true},
		{name: "empty is valid", email: "", valid: true},

		// Invalid emails - missing TLD
		{name: "no TLD", email: "pascal@laenen", valid: false, errMsg: "Invalid email format"},
		{name: "no domain", email: "test@", valid: false, errMsg: "Invalid email format"},

		// Invalid emails - format errors
		{name: "no at symbol", email: "testexample.com", valid: false, errMsg: "Invalid email format"},
		{name: "double at", email: "test@@example.com", valid: false, errMsg: "Invalid email format"},
		{name: "space in email", email: "test @example.com", valid: false, errMsg: "Invalid email format"},
		{name: "no local part", email: "@example.com", valid: false, errMsg: "Invalid email format"},
		{name: "trailing dot", email: "test@example.com.", valid: false, errMsg: "Invalid email format"},
		{name: "leading dot in domain", email: "test@.example.com", valid: false, errMsg: "Invalid email format"},
		{name: "double dot in domain", email: "test@example..com", valid: false, errMsg: "Invalid email format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Email(tt.email, tt.checkMX)

			if result.Valid != tt.valid {
				t.Errorf("Email(%q) valid = %v, want %v", tt.email, result.Valid, tt.valid)
			}

			if !tt.valid && result.Error != tt.errMsg {
				t.Errorf("Email(%q) error = %q, want %q", tt.email, result.Error, tt.errMsg)
			}
		})
	}
}

func TestEmail_Domain(t *testing.T) {
	result := Email("user@example.com", false)
	if result.Domain != "example.com" {
		t.Errorf("Email domain = %q, want %q", result.Domain, "example.com")
	}
}

func TestEmail_MXCheck(t *testing.T) {
	// Test with a known good domain (gmail.com has MX records)
	result := Email("test@gmail.com", true)
	if !result.Valid {
		t.Errorf("Email with MX check for gmail.com should be valid, got error: %s", result.Error)
	}

	// Test with a non-existent domain
	result = Email("test@thisdomain-does-not-exist-xyz123.com", true)
	if result.Valid {
		t.Error("Email with MX check for non-existent domain should be invalid")
	}
	if result.Error != "Domain does not accept email" {
		t.Errorf("Expected 'Domain does not accept email' error, got: %s", result.Error)
	}
}
