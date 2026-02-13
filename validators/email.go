package validators

import (
	"net"
	"regexp"
	"strings"
)

// emailRegex validates email format requiring a TLD (e.g., .com, .org).
// Local part: letters, digits, and special chars
// Domain: at least one dot with a valid TLD (2+ chars)
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`)

// EmailResult holds the result of email validation.
type EmailResult struct {
	Valid  bool
	Error  string
	Domain string // The domain part of the email (if valid format)
}

// Email validates an email address.
// If checkMX is true, it verifies the domain has MX records.
func Email(value string, checkMX bool) EmailResult {
	// Empty value is valid (use HTML required attribute for mandatory fields)
	if value == "" {
		return EmailResult{Valid: true}
	}

	// Format validation
	if !emailRegex.MatchString(value) {
		return EmailResult{
			Valid: false,
			Error: "Invalid email format",
		}
	}

	// Extract domain
	parts := strings.Split(value, "@")
	if len(parts) != 2 {
		return EmailResult{
			Valid: false,
			Error: "Invalid email format",
		}
	}
	domain := parts[1]

	// Optional MX record check
	if checkMX {
		mxRecords, err := net.LookupMX(domain)
		if err != nil || len(mxRecords) == 0 {
			return EmailResult{
				Valid:  false,
				Error:  "Domain does not accept email",
				Domain: domain,
			}
		}
	}

	return EmailResult{
		Valid:  true,
		Domain: domain,
	}
}
