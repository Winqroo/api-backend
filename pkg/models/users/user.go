package users

import (
	"time"
	customtypes "winqroo/pkg/customTypes"

	"github.com/google/uuid"
)

type UserProfileModel struct {
	UserType    string    `json:"userType"    dynamodbav:"userType"`
	UserID      string    `json:"userId"      dynamodbav:"userId"`
	UserEmailId string    `json:"userEmailId" dynamodbav:"userEmailId"`
	Password    string    `json:"password"    dynamodbav:"password"`
	Name        string    `json:"name"        dynamodbav:"name"`
	CreatedAt   time.Time `json:"createdAt"   dynamodbav:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"   dynamodbav:"updatedAt"`
}

func NewUserProfileModel(usr customtypes.UserSignupRequestModel) *UserProfileModel {
	now := time.Now()
	return &UserProfileModel{
		UserType: usr.UserType,
		UserID: uuid.NewString(),
		UserEmailId: usr.UserEmailId,
		Password: usr.Password,
		Name: usr.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *UserProfileModel) updateTimestamp() {
	u.UpdatedAt = time.Now()
}
