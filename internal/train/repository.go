package train

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/reww406/linetracker/config"
	"github.com/reww406/linetracker/internal/metro"
	"github.com/sirupsen/logrus"
)

type GetNextTrainsRequest struct {
	LineCode     metro.LineCode
	LocationCode string
	Direction    string
}

// Since TrainList is not a ptr, any changes will not modify the original one.
func InsertTrains(
	ctx context.Context, client *dynamodb.Client, trainList trainList,
) error {
	ddbTrains := trainList.toTrainModels()
	log.WithFields(logrus.Fields{
		"trains_len": len(ddbTrains),
	}).Info("Inserting Trains into DDB")

	for _, train := range ddbTrains {
		item, err := attributevalue.MarshalMap(train)
		if err != nil {
			return fmt.Errorf("failed to marshal train: %w", err)
		}
		log.WithFields(logrus.Fields{
			"train": train,
		}).Info("inserting train.")
		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: config.TrainTableName,
			Item:      item,
		})
		if err != nil {
			return fmt.Errorf(
				"failed to insert train with location code: %s created: %d with error: %w",
				train.LocationCode, train.CreatedEpochMs, err,
			)
		}
	}
	return nil
}

// Line -> Location -> Direction
func GetTrainPredictions(
	ctx context.Context, client *dynamodb.Client, request GetNextTrainsRequest,
) ([]TrainModel, error) {
	validMinutes := -10 * time.Minute
	timeRange := time.Now().Add(validMinutes).UnixMilli()

	keyExpr := expression.Key("locationCode").
		Equal(expression.Value(request.LocationCode)).
		And(expression.Key("createdEpochMs").
			GreaterThanEqual(expression.Value(timeRange)))

	filterExpr := expression.Name("lineCode").
		Equal(expression.Value(request.LineCode)).
		And(expression.Name("destination").
			Equal(expression.Value(request.Direction)))

	expr, err := expression.NewBuilder().
		WithKeyCondition(keyExpr).
		WithFilter(filterExpr).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build ddb expression %w", err)
	}

	// Perform the query
	result, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 config.TrainTableName,
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query trains: %w", err)
	}

	log.WithField("result_len", len(result.Items)).Info("trains found.")

	// Convert results to TrainModels
	trains := make([]TrainModel, len(result.Items))
	for i, item := range result.Items {
		trains[i] = itemToDdbTrain(item)
	}

	return trains, nil
}

func itemToDdbTrain(item map[string]types.AttributeValue) TrainModel {
	carCount, _ := strconv.Atoi(item["carCount"].(*types.AttributeValueMemberN).Value)
	minutes, _ := strconv.Atoi(item["minutes"].(*types.AttributeValueMemberN).Value)
	createdEpochMs, _ := strconv.Atoi(item["createdEpochMs"].(*types.AttributeValueMemberN).Value)
	return TrainModel{
		CarCount:        int8(carCount),
		Destination:     item["destination"].(*types.AttributeValueMemberS).Value,
		DestinationCode: item["destinationCode"].(*types.AttributeValueMemberS).Value,
		DestinationName: item["destinationName"].(*types.AttributeValueMemberS).Value,
		Group:           item["group"].(*types.AttributeValueMemberS).Value,
		LineCode:        item["lineCode"].(*types.AttributeValueMemberS).Value,
		LocationCode:    item["locationCode"].(*types.AttributeValueMemberS).Value,
		LocationName:    item["locationName"].(*types.AttributeValueMemberS).Value,
		Minutes:         int8(minutes),
		CreatedEpochMs:  int64(createdEpochMs),
	}
}
