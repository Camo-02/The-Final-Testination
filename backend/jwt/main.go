package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// FIXME: should be an env variable
const JWT_SECRET = "secret"

type FinalTestinationClaims struct {
	jwt.RegisteredClaims
	PlayerID string `json:"user_id"`
}

func GenerateJWT(playerID string) (string, error) {
	claims := FinalTestinationClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "the-final-testination",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		PlayerID: playerID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	return tokenString, err
}

func ParseJWT(tokenString string) (*FinalTestinationClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &FinalTestinationClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*FinalTestinationClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("invalid JWT")
	}
}

func ParseJWTWithClaims(tokenString string) (*FinalTestinationClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &FinalTestinationClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*FinalTestinationClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("couldn't parse claims")
	}
}
