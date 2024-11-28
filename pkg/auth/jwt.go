package jwt

import (
	"errors"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
)

const UserClaimKey = "User-Claims"

func CreateToken(secret []byte, claims *UserClaims) (string, error) {
	claims.Role = "user"
	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(secret)
}

func ParseToken(tokenString string, secret []byte) (*UserClaims, error) {
	if strings.Count(tokenString, ".") != 2 {
		return nil, errors.New("token contains an invalid number of segments")
	}
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
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
