package users

import (
	"context"
	"fmt"

	"winqroo/config"
	userModels "winqroo/pkg/models/users"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserProfileRepo struct {
	Client *dynamodb.Client
}

func NewUserProfileRepo(c *dynamodb.Client) *UserProfileRepo {
	return &UserProfileRepo{Client: c}
}

func (r *UserProfileRepo) GetUserProfileInfoByID(
	ctx context.Context,
	userType string,
	userID string,
) (userProfile *userModels.UserProfileModel, err error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		userType: userType,
		userID:   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to dynamoAttr marshal map: %w", err)
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(config.GetUserProfilesStore()),
		Key:       key,
	}
	output, err := r.Client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get userProfile from db: %w", err)
	}
	// Check if the Item field in the response is nil or empty
	if output.Item == nil || len(output.Item) == 0 {
		return nil, nil
	}

	err = attributevalue.UnmarshalMap(output.Item, &userProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to dynamoAttr marshal map: %w", err)
	}

	return userProfile, nil
}

func (r *UserProfileRepo) GetUserProfileInfoByEmail(
	ctx context.Context,
	userType string,
	userEmail string,
) (userProfile *userModels.UserProfileModel, err error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(config.GetUserProfilesStore()),
		IndexName:              aws.String(config.GetUserProfilesStoreEmailIndex()),
		KeyConditionExpression: aws.String("userType = :userType AND userEmail = :userEmail"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userType":  &types.AttributeValueMemberS{Value: userType},
			":userEmail": &types.AttributeValueMemberS{Value: userEmail},
		},
	}
	output, err := r.Client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get userProfile from db: %w", err)
	}
	// Check if the Item field in the response is nil or empty
	if output.Items == nil || len(output.Items) == 0 {
		return nil, nil
	}

	err = attributevalue.UnmarshalMap(output.Items[0], &userProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to dynamoAttr marshal map: %w", err)
	}

	return userProfile, nil
}

func (r *UserProfileRepo) PutUserProfileInfo(
	ctx context.Context,
	user *userModels.UserProfileModel,
) (err error) {
	data, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("failed to parse condition map: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(config.GetUserProfilesStore()),
		Item:      data,
	}

	_, err = r.Client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put user profile info: %w", err)
	}

	return nil
}

func (r *UserProfileRepo) DeleteUserProfileInfoByID(
	ctx context.Context,
	userType string,
	userID string,
) (err error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		userType: userType,
		userID:   userID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal key: %w", err)
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(config.GetUserProfilesStore()),
		Key:       key,
	}

	_, err = r.Client.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete user profile: %w", err)
	}

	return nil
}
