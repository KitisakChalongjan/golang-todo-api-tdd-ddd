package helper

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessTokenWithClaims(claims jwt.MapClaims, secretKey string) (string, error) {

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("fail to generate accessToken: %w", err)
	}

	return accessTokenString, nil
}

func TokenClaimsFromAccessTokenString(jwtString string, secretKey string) (jwt.Token, error) {

	token, err := jwt.ParseWithClaims(
		jwtString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {

			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		return jwt.Token{}, err
	}

	return *token, nil
}
