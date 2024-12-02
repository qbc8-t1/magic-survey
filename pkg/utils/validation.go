package utils

import (
	"regexp"
	"time"
)

// IsValidEmail checks if an mail address is in a valid format.
func IsValidEmail(email string) bool {
	// Simple mail regex for validation
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

// IsAllDigits checks if a string contains only numeric digits.
func IsAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func IsValidDate(dateStr string) bool {
	_, err := time.Parse("2006-01-02 15:04:05", dateStr)
	return err == nil
}
