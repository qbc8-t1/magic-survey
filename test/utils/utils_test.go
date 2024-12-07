package test

import (
	"testing"

	"github.com/QBC8-Team1/magic-survey/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestIsValidEmail(t *testing.T) {
	// Valid email addresses
	t.Run("ValidEmail", func(t *testing.T) {
		isValid := utils.IsValidEmail("foo@example.com")
		assert.True(t, isValid, "expected is valid to be true, but got %v", isValid)
	})

	t.Run("ValidEmailWithSubdomain", func(t *testing.T) {
		isValid := utils.IsValidEmail("foo@sub.example.com")
		assert.True(t, isValid, "expected is valid to be true, but got %v", isValid)
	})

	t.Run("ValidEmailWithNumbers", func(t *testing.T) {
		isValid := utils.IsValidEmail("user123@example123.com")
		assert.True(t, isValid, "expected is valid to be true, but got %v", isValid)
	})

	t.Run("ValidEmailWithSpecialCharacters", func(t *testing.T) {
		isValid := utils.IsValidEmail("foo.bar+123@example.com")
		assert.True(t, isValid, "expected is valid to be true, but got %v", isValid)
	})

	// Invalid email addresses
	t.Run("MissingAtSymbol", func(t *testing.T) {
		isValid := utils.IsValidEmail("fooexample.com")
		assert.False(t, isValid, "expected is valid to be false, but got %v", isValid)
	})

	t.Run("MissingDomain", func(t *testing.T) {
		isValid := utils.IsValidEmail("foo@")
		assert.False(t, isValid, "expected is valid to be false, but got %v", isValid)
	})

	t.Run("MissingUsername", func(t *testing.T) {
		isValid := utils.IsValidEmail("@example.com")
		assert.False(t, isValid, "expected is valid to be false, but got %v", isValid)
	})

	t.Run("InvalidCharacters", func(t *testing.T) {
		isValid := utils.IsValidEmail("foo@exa!mple.com")
		assert.False(t, isValid, "expected is valid to be false, but got %v", isValid)
	})

	t.Run("WhitespaceInEmail", func(t *testing.T) {
		isValid := utils.IsValidEmail("foo bar@example.com")
		assert.False(t, isValid, "expected is valid to be false, but got %v", isValid)
	})

	t.Run("TrailingDotInDomain", func(t *testing.T) {
		isValid := utils.IsValidEmail("foo@example.com.")
		assert.False(t, isValid, "expected is valid to be false, but got %v", isValid)
	})

	t.Run("InvalidTLD", func(t *testing.T) {
		isValid := utils.IsValidEmail("foo@example.c")
		assert.False(t, isValid, "expected is valid to be false, but got %v", isValid)
	})

	// Edge cases
	t.Run("EmptyString", func(t *testing.T) {
		isValid := utils.IsValidEmail("")
		assert.False(t, isValid, "expected is valid to be false, but got %v", isValid)
	})

	t.Run("OnlyWhitespace", func(t *testing.T) {
		isValid := utils.IsValidEmail("   ")
		assert.False(t, isValid, "expected is valid to be false, but got %v", isValid)
	})

	t.Run("EmailWithUnicode", func(t *testing.T) {
		isValid := utils.IsValidEmail("foo@例子.com")
		assert.False(t, isValid, "expected is valid to be false, but got %v", isValid)
	})
}

func TestIsValidNationalCode(t *testing.T) {
	tests := []struct {
		code     string
		expected bool
	}{
		{"1234567891", true},
		{"9876543210", true},
		{"12345abc890", false},
		{"123", false},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			result, err := utils.IsValidNationalCode(tt.code)
			if err != nil {
				t.Fatalf("IsValidNationalCode(%q) returned error: %v", tt.code, err)
			}
			if result != tt.expected {
				t.Errorf("IsValidNationalCode(%q) = %v; want %v", tt.code, result, tt.expected)
			}
		})
	}
}
