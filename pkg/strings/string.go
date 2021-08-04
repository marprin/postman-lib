package strings

import "regexp"

var (
	emailRegex       = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	alphaNumericOnly = regexp.MustCompile("[a-zA-Z0-9]+")
)

// IsValidEmail return true if email address has a valid format
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func AlphaNumericOnly(payload string) string {
	return alphaNumericOnly.FindString(payload)
}
