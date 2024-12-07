package utils

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

// GenerateRandomCode generates a random 6-digit code
func GenerateRandomCode() string {
	rand.Seed(uint64(time.Now().UnixNano()))
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
