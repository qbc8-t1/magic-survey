package utils

import (
	"regexp"
	"strings"
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

func IsValidBirthdate(input string) (bool, string) {
	layout := "2006-01-02"
	parsedDate, err := time.Parse(layout, input)
	if err != nil {
		return false, "Invalid date format, expected yyyy-mm-dd"
	}

	currentDate := time.Now()
	if parsedDate.After(currentDate) {
		return false, "Birthdate cannot be in the future"
	}

	formattedDate := parsedDate.Format(layout)
	if input != formattedDate {
		return false, "Date is not in the expected format yyyy-mm-dd"
	}

	return true, "Valid date"
}

func IsValidCity(cityName string) (bool, string) {
	cityName = strings.TrimSpace(cityName)
	if len(cityName) < 2 || len(cityName) > 50 {
		return false, "City name must be between 2 and 50 characters."
	}
	isValid := regexp.MustCompile("^[a-zA-Z]+$").MatchString(cityName)
	if !isValid {
		return false, "City name can only contain English letters (no digits or special characters)."
	}
	return true, "City name is valid."
}

func HasTimePassed(createdAt time.Time, seconds int) bool {
	currentTime := time.Now()
	duration := currentTime.Sub(createdAt)
	if duration >= time.Duration(seconds)*time.Second {
		return true
	}
	return false
}
