package customtypes

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserClaim struct {
	UserID string
	userType   string
	jwt.RegisteredClaims
}

func NewUserClaim(userID,userType string) *UserClaim {
	now := time.Now()
	expTime := now.Add(24 * time.Hour)
	return &UserClaim{
		UserID: userID,
		userType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "winqroo.com",
			Subject: "Auth Token",
			IssuedAt: &jwt.NumericDate{Time: now},
			ExpiresAt: &jwt.NumericDate{Time: expTime},
		},
	}
}