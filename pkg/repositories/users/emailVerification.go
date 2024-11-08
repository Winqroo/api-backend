package users

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"winqroo/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	sestypes "github.com/aws/aws-sdk-go-v2/service/ses/types"
	"golang.org/x/exp/rand"
)

type UserEmailVerificationRepo struct {
	ClientDB  *dynamodb.Client
	ClientSES *ses.Client
}

func NewUserEmailVerificationRepo(ddb *dynamodb.Client, ses *ses.Client) *UserEmailVerificationRepo {
	return &UserEmailVerificationRepo{
		ClientDB:  ddb,
		ClientSES: ses,
	}
}

func generateOtpAndAuthCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func sendEmail(ctx context.Context, client *ses.Client, recipient string, otp string) error {
	senderEmail := config.GetSESOtpSenderEmail()

	subject := "Email-verification"
	body := fmt.Sprintf("Your OTP for email verification is: %s. Please enter this OTP to complete your verification process.", otp)

	input := &ses.SendEmailInput{
		Source: aws.String(senderEmail),
		Destination: &sestypes.Destination{
			ToAddresses: []string{recipient},
		},
		Message: &sestypes.Message{
			Subject: &sestypes.Content{
				Data: aws.String(subject),
			},
			Body: &sestypes.Body{
				Text: &sestypes.Content{
					Data: aws.String(body),
				},
			},
		},
	}

	_, err := client.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (r *UserEmailVerificationRepo) GetOtpToRegister(ctx context.Context, email string) (string, error) {
	otp := generateOtpAndAuthCode()
	authCode := generateOtpAndAuthCode()

	err := sendEmail(ctx, r.ClientSES, email, otp)
	if err != nil {
		return "", fmt.Errorf("failed to send OTP email: %w", err)
	}

	ttl := time.Now().Add(5 * time.Minute).Unix() // OTP valid for 5 minutes
	item := map[string]types.AttributeValue{
		"Email":    &types.AttributeValueMemberS{Value: email},
		"OTP":      &types.AttributeValueMemberS{Value: otp},
		"AuthCode": &types.AttributeValueMemberS{Value: authCode},
		"TTL":      &types.AttributeValueMemberN{Value: strconv.FormatInt(ttl, 10)},
	}

	_, err = r.ClientDB.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(config.GetUserOtpStore()),
		Item:      item,
	})
	if err != nil {
		return "", fmt.Errorf("failed to store OTP in DynamoDB: %w", err)
	}

	return authCode, nil
}

func (r *UserEmailVerificationRepo) ResendOtp(ctx context.Context, email string) error {
	result, err := r.ClientDB.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(config.GetUserOtpStore()),
		Key: map[string]types.AttributeValue{
			"Email": &types.AttributeValueMemberS{Value: email},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to retrieve email from DynamoDB: %w", err)
	}
	if result.Item == nil {
		return fmt.Errorf("no OTP request found for email: %s", email)
	}

	otp := generateOtpAndAuthCode()
	err = sendEmail(ctx, r.ClientSES, email, otp)
	if err != nil {
		return fmt.Errorf("failed to resend OTP email: %w", err)
	}

	ttl := time.Now().Add(5 * time.Minute).Unix()
	_, err = r.ClientDB.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.GetUserOtpStore()),
		Key: map[string]types.AttributeValue{
			"Email": &types.AttributeValueMemberS{Value: email},
		},
		ExpressionAttributeNames: map[string]string{
			"#O": "OTP",
			"#T": "TTL",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":o": &types.AttributeValueMemberS{Value: otp},
			":t": &types.AttributeValueMemberN{Value: strconv.FormatInt(ttl, 10)},
		},
		UpdateExpression: aws.String("SET #O = :o, #T = :t"),
	})
	if err != nil {
		return fmt.Errorf("failed to update OTP in DynamoDB: %w", err)
	}

	return nil
}

func (r *UserEmailVerificationRepo) VerifyOtp(ctx context.Context, email, inputOtp, AuthCode string) (bool, error) {
	result, err := r.ClientDB.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(config.GetUserOtpStore()),
		Key: map[string]types.AttributeValue{
			"Email": &types.AttributeValueMemberS{Value: email},
		},
	})
	if err != nil {
		return false, fmt.Errorf("failed to retrieve OTP from DynamoDB: %w", err)
	}

	if result.Item == nil {
		return false, fmt.Errorf("no OTP found for email: %s", email)
	}

	var dbItem struct {
		Email    string
		OTP      string
		AuthCode string
		TTL      int64
	}
	err = attributevalue.UnmarshalMap(result.Item, &dbItem)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal OTP data: %w", err)
	}

	currentTime := time.Now().Unix()
	if dbItem.OTP != inputOtp {
		return false, fmt.Errorf("incorrect OTP")
	}
	if dbItem.TTL < currentTime {
		return false, fmt.Errorf("OTP has expired")
	}

	return true, nil
}
