package user

import (
	"context"
	"winqroo/pkg/utils"
)

func (u *UserVerificationService) GetOtpToRegister(ctx context.Context, email string) (interface{}, *utils.CustomError) {
	authCode, err := u.emailVerificationRepo.GetOtpToRegister(ctx,email)
	if err != nil {
		return nil, utils.ReturnError(err,utils.ErrCodes.Common.ErrCodeInternalServerError)
	}

	return struct {
		email string
		authCode string
	}{
		email: email,
		authCode: authCode,
	}, nil
}

func (u *UserVerificationService) ResendOtp(ctx context.Context, email string) (interface{}, *utils.CustomError) {
	err := u.emailVerificationRepo.ResendOtp(ctx,email)
	if err != nil {
		return nil, utils.ReturnError(err,utils.ErrCodes.Common.ErrCodeInternalServerError)
	}

	return struct {
		email string
	}{
		email: email,
	}, nil
}

func (u *UserVerificationService) VerifyOtp(ctx context.Context, email, otp, authCode string) (interface{}, *utils.CustomError) {
	success, err := u.emailVerificationRepo.VerifyOtp(ctx,email,otp,authCode)
	if err != nil {
		return nil, utils.ReturnError(err,utils.ErrCodes.Common.ErrCodeInternalServerError)
	}

	return struct {
		success bool
	}{
		success: success,
	}, nil
}
