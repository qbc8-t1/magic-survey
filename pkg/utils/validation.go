package utils

import (
	"log"
	"regexp"
	"strconv"
	"strings"
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

func IsValidNationalCode(code string) (bool, error) {
	reg, err := regexp.Compile("/[^0-9]/")

	if err != nil {
		log.Fatal(err)
	}

	code = reg.ReplaceAllString(code, "")

	if len(code) != 10 {
		return false, nil
	}

	codes := strings.Split(code, "")

	last, err := strconv.Atoi(codes[9])
	if err != nil {
		return false, err
	}

	i := 10
	sum := 0

	for in, el := range codes {
		temp, err := strconv.Atoi(el)

		if err != nil {
			log.Fatal(err)
		}

		if in == 9 {
			break
		}

		sum += temp * i
		i -= 1
	}

	mod := sum % 11

	if mod >= 2 {
		mod = 11 - mod
	}

	return mod == last, nil
}
