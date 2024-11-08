package tasks

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type TaskProgressRepo struct {
	Client *dynamodb.Client
}

func NewTaskProgressRepo(c *dynamodb.Client) *TaskProgressRepo {
	return &TaskProgressRepo{Client: c}
}

//uid(pk) to taskid(sk)
//taskid to uid
//uid to createdat
//taskid to createdat

