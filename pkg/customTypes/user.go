package customtypes

type UserLoginRequestModel struct {
	UserType    string `json:"userType"    dynamodbav:"userType"`
	UserEmailId string `json:"userEmailId" dynamodbav:"userEmailId"`
	Password    string `json:"password"    dynamodbav:"password"`
}

type UserSignupRequestModel struct {
	UserType    string `json:"userType"    dynamodbav:"userType"`
	UserEmailId string `json:"userEmailId" dynamodbav:"userEmailId"`
	Password    string `json:"password"    dynamodbav:"password"`
	Name        string `json:"name"        dynamodbav:"name"`
}
