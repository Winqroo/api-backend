package tasks

import (
	"context"
	"fmt"

	"winqroo/config"
	taskModels "winqroo/pkg/models/tasks"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TaskInfoRepo struct {
	Client *dynamodb.Client
}

func NewTaskInfoRepo(c *dynamodb.Client) *TaskInfoRepo {
	return &TaskInfoRepo{Client: c}
}

func (r *TaskInfoRepo) GetTaskInfoByID(
	ctx context.Context,
	taskID string,
) (task *taskModels.TaskInfoModel, err error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		taskID: taskID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to dynamoAttr marshal map: %w", err)
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(config.GetTaskInfoStore()),
		Key:       key,
	}
	output, err := r.Client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get Task from db: %w", err)
	}
	// Check if the Item field in the response is nil or empty
	if output.Item == nil || len(output.Item) == 0 {
		return nil, nil
	}

	err = attributevalue.UnmarshalMap(output.Item, &task)
	if err != nil {
		return nil, fmt.Errorf("failed to dynamoAttr marshal map: %w", err)
	}

	return task, nil
}

func (r *TaskInfoRepo) PutTaskInfo(
	ctx context.Context,
	task *taskModels.TaskInfoModel,
) (err error) {
	data, err := attributevalue.MarshalMap(task)
	if err != nil {
		return fmt.Errorf("failed to parse condition map: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(config.GetTaskInfoStore()),
		Item:      data,
	}

	_, err = r.Client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put task: %w", err)
	}

	return nil
}

func (r *TaskInfoRepo) DeleteTaskInfoByID(
	ctx context.Context,
	taskID string,
) (err error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		taskID: taskID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal key: %w", err)
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(config.GetTaskInfoStore()),
		Key:       key,
	}

	_, err = r.Client.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

func (r *TaskInfoRepo) GetAllTaskInfo(
	ctx context.Context,
	taskStatus taskModels.TaskInfoStatus,
	limit int,
	lastEvaluatedKey interface{},
) (result []*taskModels.TaskInfoModel, lastKey map[string]types.AttributeValue, err error) {
	expAttrVals := map[string]types.AttributeValue{}

	filterExpression := ""
	if taskStatus != "" {
		expAttrVals[":taskStatus"] = &types.AttributeValueMemberS{Value: string(taskStatus)}
		filterExpression = "taskStatus = :taskStatus"
	}

	queryInput := &dynamodb.ScanInput{
		TableName:                 aws.String(config.GetTaskInfoStore()),
		ExpressionAttributeValues: expAttrVals,
		Limit:                     aws.Int32(int32(limit)),
	}

	if lastEvaluatedKey != nil {
		if lek, ok := lastEvaluatedKey.(map[string]types.AttributeValue); ok {
			queryInput.ExclusiveStartKey = lek
		}
	}

	if filterExpression != "" {
		queryInput.FilterExpression = aws.String(filterExpression)
	}

	for {
		resp, err := r.Client.Scan(ctx, queryInput)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan tasks info: %w", err)
		}

		var transactions []*taskModels.TaskInfoModel
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &transactions); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal tasks items: %w", err)
		}

		result = append(result, transactions...)

		if len(result) >= limit {
			result = result[:limit]
			lastKey = resp.LastEvaluatedKey
			break
		}

		if resp.LastEvaluatedKey == nil {
			break
		}

		queryInput.ExclusiveStartKey = resp.LastEvaluatedKey
	}

	return result, lastKey, nil
}
