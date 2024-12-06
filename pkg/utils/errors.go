package utils

import "errors"

// ErrorIncludes a helper function to check if the error is in the known errors array
func ErrorIncludes(err error, knownErrors []error) bool {
	for _, knownError := range knownErrors {
		if errors.Is(err, knownError) {
			return true
		}
	}
	return false
}
