package train

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/reww406/linetracker/config"
	"github.com/sirupsen/logrus"
)

// Since TrainList is not a ptr, any changes will not modify the original one.
func InsertTrains(
	ctx context.Context, client *dynamodb.Client, trainList TrainList,
) error {
	ddbTrains := trainList.toDdbTrains()
	// For some reason the first 100 or so trains are empty...
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
			return fmt.Errorf("failed to insert train %s: %w", train.Id, err)
		}
	}
	return nil
}

func itemToDdbTrain(item map[string]types.AttributeValue) DdbTrain {
	carCount, _ := strconv.Atoi(item["carCount"].(*types.AttributeValueMemberN).Value)
	minutes, _ := strconv.Atoi(item["minutes"].(*types.AttributeValueMemberN).Value)
	createdEpoch, _ := strconv.Atoi(item["createdEpoch"].(*types.AttributeValueMemberN).Value)
	return DdbTrain{
		CarCount:        int8(carCount),
		Destination:     item["destination"].(*types.AttributeValueMemberS).Value,
		DestinationCode: item["destinationCode"].(*types.AttributeValueMemberS).Value,
		DestinationName: item["destinationName"].(*types.AttributeValueMemberS).Value,
		Group:           item["group"].(*types.AttributeValueMemberS).Value,
		LineCode:        item["lineCode"].(*types.AttributeValueMemberS).Value,
		LocationCode:    item["locationCode"].(*types.AttributeValueMemberS).Value,
		LocationName:    item["locationName"].(*types.AttributeValueMemberS).Value,
		Minutes:         int8(minutes),
		CreatedEpoch:    int64(createdEpoch),
		Id:              item["id"].(*types.AttributeValueMemberS).Value,
	}
}



