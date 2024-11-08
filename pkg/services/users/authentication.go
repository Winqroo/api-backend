package user

import (
	"context"
	"fmt"

	customtypes "winqroo/pkg/customTypes"
	userModels "winqroo/pkg/models/users"
	"winqroo/pkg/utils"
)

func (u *UserAuthenticationService) LoginUser(ctx context.Context, userType, userEmail, password string) (string, *utils.CustomError) {
	user, err := u.userProfileRepo.GetUserProfileInfoByEmail(ctx, userType, userEmail)
	if err != nil {
		return "", utils.ReturnError(err, utils.ErrCodes.Common.ErrCodeInternalServerError)
	}
	if user == nil {
		return "", utils.ReturnError(fmt.Errorf("User Not Found"), utils.ErrCodes.Common.ErrCodeRecordNotFound)
	}

	isValidPassword, err := u.hashingSystem.ComparePassword(user.Password, password)
	if err != nil {
		return "", utils.ReturnError(err, utils.ErrCodes.Common.ErrCodeInternalServerError)
	}
	if !isValidPassword {
		return "", utils.ReturnError(fmt.Errorf("Invalid Password"), utils.ErrCodes.Common.ErrUnAuthorised)
	}

	userClaim := customtypes.NewUserClaim(user.UserID, user.UserType)
	jwtTokenString, err := u.jwtAuthSystem.GenerateSessionJWT(userClaim)
	if err != nil {
		return "", utils.ReturnError(err, utils.ErrCodes.Common.ErrCodeInternalServerError)
	}

	return jwtTokenString, nil
}

func (u *UserAuthenticationService) RegisterNewUser(ctx context.Context, reqUser customtypes.UserSignupRequestModel) *utils.CustomError {
	user, err := u.userProfileRepo.GetUserProfileInfoByEmail(ctx, reqUser.UserType, reqUser.UserEmailId)
	if err != nil {
		return utils.ReturnError(err, utils.ErrCodes.Common.ErrCodeInternalServerError)
	}
	if user != nil {
		return utils.ReturnError(fmt.Errorf("User with %s already exists",reqUser.UserEmailId),utils.ErrCodes.Common.ErrCodeInternalServerError)
	}

	reqUser.Password, err = u.hashingSystem.HashPassword(reqUser.Password)
	if err != nil {
		return utils.ReturnError(err, utils.ErrCodes.Common.ErrCodeInternalServerError)
	}

	user = userModels.NewUserProfileModel(reqUser)
	err = u.userProfileRepo.PutUserProfileInfo(ctx, user)
	if err != nil {
		return utils.ReturnError(err,utils.ErrCodes.Common.ErrCodeInternalServerError)
	}

	return nil
}
