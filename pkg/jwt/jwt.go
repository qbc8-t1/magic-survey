package jwt

import (
	"errors"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"strings"
)

func CreateToken(secret []byte, claims *UserClaims) (string, error) {
	return jwt2.NewWithClaims(jwt2.SigningMethodHS512, claims).SignedString(secret)
}

// ParseToken validates the JWT token and returns the claims
func ParseToken(tokenString string, secret []byte) (*UserClaims, error) {
	if strings.Count(tokenString, ".") != 2 {
		return nil, errors.New("token contains an invalid number of segments")
	}
	token, err := jwt2.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt2.Token) (interface{}, error) {
		return secret, nil
	})

	var claim *UserClaims
	if token.Claims != nil {
		cc, ok := token.Claims.(*UserClaims)

		if ok {
			claim = cc
		}
	}

	if err != nil {
		return claim, err
	}

	if !token.Valid {
		return claim, errors.New("token is not valid")
	}

	return claim, nil
}
