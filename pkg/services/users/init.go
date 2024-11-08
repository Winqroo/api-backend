package user

import (
	"context"

	userModels "winqroo/pkg/models/users"
	"winqroo/pkg/utils"
)

type (
	UserProfileRepo interface {
		GetUserProfileInfoByID(
			ctx context.Context,
			userType string,
			userID string,
		) (userProfile *userModels.UserProfileModel, err error)
		PutUserProfileInfo(
			ctx context.Context,
			user *userModels.UserProfileModel,
		) (err error)
		DeleteUserProfileInfoByID(
			ctx context.Context,
			userType string,
			userID string,
		) (err error)
		GetUserProfileInfoByEmail(
			ctx context.Context,
			userType string,
			userEmail string,
		) (userProfile *userModels.UserProfileModel, err error)
	}

	EmailVerificationRepo interface {
		GetOtpToRegister(
			ctx context.Context,
			email string,
		) (string, error)
		ResendOtp(
			ctx context.Context,
			email string,
		) error
		VerifyOtp(
			ctx context.Context,
			email string,
			inputOtp string,
			AuthCode string,
		) (bool, error)
	}
)

type (
	UserAuthenticationService struct {
		userProfileRepo UserProfileRepo
		hashingSystem   *utils.HashingSystemUtils
		jwtAuthSystem   *utils.JwtAuthSystemUtils
	}
	UserVerificationService struct {
		userProfileRepo       UserProfileRepo
		emailVerificationRepo EmailVerificationRepo
	}
	UserInteractionsService struct {
		userProfileRepo UserProfileRepo
	}
)

func NewUserVerificationService(
	emailVerificationRepo EmailVerificationRepo,
) *UserVerificationService {
	return &UserVerificationService{
		emailVerificationRepo: emailVerificationRepo,
	}
}

func NewUserInteractionsService(
	userProfileRepo UserProfileRepo,
) *UserInteractionsService {
	return &UserInteractionsService{
		userProfileRepo: userProfileRepo,
	}
}

func NewUserAuthenticationService(
	userProfileRepo UserProfileRepo,
) *UserAuthenticationService {
	return &UserAuthenticationService{
		userProfileRepo: userProfileRepo,
		hashingSystem:   utils.NewHashingSystemUtils(),
		jwtAuthSystem:   utils.NewJwtAuthSystemUtils(),
	}
}
