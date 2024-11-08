package utils

import (
	"fmt"
	"winqroo/config"
	"winqroo/pkg/customTypes"

	"github.com/golang-jwt/jwt/v4"
)

type JwtAuthSystemUtils struct {
	SecretKey string
}

func NewJwtAuthSystemUtils() *JwtAuthSystemUtils {
	return &JwtAuthSystemUtils{
		SecretKey: config.GetJWTAuthSecretKey(),
	}
}

func (j *JwtAuthSystemUtils) GenerateSessionJWT(claims *customtypes.UserClaim) (tokenString string, err error) {
		if claims.Valid() != nil {
			return "", fmt.Errorf("Invalid claims")
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err = token.SignedString([]byte(j.SecretKey))
		if err != nil {
			return "", err
		}

		return tokenString, nil
}

func (j *JwtAuthSystemUtils) ValidateSessionJWT(tokenString string, claims *customtypes.UserClaim) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}