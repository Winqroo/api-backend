package config

import "os"

type AwsConfig struct {
	Region    string
	AccessKey string
	SecretKey string
}

type Config struct {
	Env                         string
	AwsConfig                   AwsConfig
	DynamoEndpoint              string
	ServerPort                  string
	UserProfilesStore           string
	UserProfilesStoreEmailIndex string
	TaskInfoStore               string
	UserOtpStore string
	OTPManagementStore          string
	JWTAuthSecretKey            string
	HashingSecretKey            string
	SESOtpSenderEmail              string
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

var instance *Config

func Init() {
	if instance == nil {
		instance = &Config{
			Env: getEnvWithDefault("ENV", "dev"),
			AwsConfig: AwsConfig{
				Region:    getEnvWithDefault("AWS_REGION", ""),
				AccessKey: getEnvWithDefault("AWS_ACCESS_KEY_ID", ""),
				SecretKey: getEnvWithDefault("AWS_SECRET_ACCESS_KEY", ""),
			},
			ServerPort:                  getEnvWithDefault("SERVER_PORT", "3000"),
			DynamoEndpoint:              getEnvWithDefault("DYNAMO_ENDPOINT", ""),
			UserProfilesStore:           getEnvWithDefault("USER_PROFILES_STORE", ""),
			UserProfilesStoreEmailIndex: getEnvWithDefault("USER_PROFILES_STORE_EMAIL_INDEX", ""),
			TaskInfoStore:               getEnvWithDefault("TASK_INFO_STORE", ""),
			UserOtpStore:                getEnvWithDefault("USER_OTP_STORE", ""),
			OTPManagementStore:          getEnvWithDefault("OTP_MANAGEMENT_STORE", ""),
			JWTAuthSecretKey:            getEnvWithDefault("JWT_AUTH_SECRET_KEY", ""),
			HashingSecretKey:            getEnvWithDefault("HASHING_SECRET_KEY", ""),
			SESOtpSenderEmail:              getEnvWithDefault("SES_SENDER_EMAIL", ""),
		}
	}
}

func GetInstance() *Config {
	return instance
}

func GetUserProfilesStore() string {
	return instance.UserProfilesStore
}

func GetUserProfilesStoreEmailIndex() string {
	return instance.UserProfilesStoreEmailIndex
}

func GetTaskInfoStore() string {
	return instance.TaskInfoStore
}

func GetUserOtpStore() string {
	return instance.UserOtpStore
}

func GetJWTAuthSecretKey() string {
	return instance.JWTAuthSecretKey
}

func GetHashingSecretKey() string {
	return instance.HashingSecretKey
}

func GetSESOtpSenderEmail() string {
	return instance.SESOtpSenderEmail
}
