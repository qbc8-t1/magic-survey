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
