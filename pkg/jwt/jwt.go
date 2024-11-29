package jwt

import (
	"time"

	jwt2 "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("test")

func GenerateToken(userID uint) (string, error) {
	claims := jwt2.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expiration (24 hours)
		"iat":     time.Now().Unix(),
	}

	token := jwt2.NewWithClaims(jwt2.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

// ValidateToken validates the JWT token and returns the claims
func ValidateToken(tokenString string) (jwt2.MapClaims, error) {
	token, err := jwt2.Parse(tokenString, func(token *jwt2.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt2.SigningMethodHMAC); !ok {
			return nil, jwt2.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt2.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt2.ErrTokenSignatureInvalid
}

// GetTokenExpiry returns the token expiry time
func GetTokenExpiry() time.Time {
	return time.Now().Add(24 * time.Hour)
}
