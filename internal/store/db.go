package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	appConfig "github.com/reww406/linetracker/config"
	"github.com/reww406/linetracker/internal/station"
	"github.com/sirupsen/logrus"
)

var log = appConfig.GetLogger()

func tableExists(ctx context.Context, client *dynamodb.Client, tableName *string) bool {
	_, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: tableName,
	})
	return err == nil
}

func tableItemCount(
	ctx context.Context, client *dynamodb.Client, tableName *string,
) (int, error) {
	input := &dynamodb.DescribeTableInput{
		TableName: tableName,
	}

	result, err := client.DescribeTable(ctx, input)
	if err != nil {
		return 0, err
	}
	// De-reference pointer
	return int(*result.Table.ItemCount), nil
}

func createStationsTable(client *dynamodb.Client) error {
	log.WithFields(logrus.Fields{
		"TableName": appConfig.StationTableName,
	}).Info("Creating DDB Table")
	_, err := client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName: appConfig.StationTableName,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("code"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("code"),
				KeyType:       types.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})
	if err != nil {
		return fmt.Errorf("error creating stations table: %w", err)
	}
	return nil
}

func createTrainTable(client *dynamodb.Client) error {
	log.WithFields(logrus.Fields{
		"TableName": appConfig.TrainTableName,
	}).Info("Creating DDB Table")
	_, err := client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName: appConfig.TrainTableName,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("createdEpoch"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("createdEpoch"),
				KeyType:       types.KeyTypeRange,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})
	if err != nil {
		return fmt.Errorf("error creating stations table: %w", err)
	}
	return nil
}

func insertStations(client *dynamodb.Client) error {
	// Insert stations if they don't already exist
	count, err := tableItemCount(context.Background(), client, appConfig.StationTableName)
	if err != nil {
		log.WithFields(logrus.Fields{
			"TableName": appConfig.StationTableName,
		}).Warn("Failed to get item count from table.")
	}
	if count <= 0 {
		log.WithFields(logrus.Fields{
			"TableName": appConfig.StationTableName,
		}).Info("inserting stations into DDB.")
		stations, err := station.GetStations()
		if err != nil {
			return fmt.Errorf("failed to get stations: %w", err) 
		}
		if err := station.InsertStations(context.Background(), client, *stations); err != nil {
			return fmt.Errorf("failed to insert stations: %w", err) 
		}
	}
	return nil
}

func InitDB() (*dynamodb.Client, error) {
	// Configure AWS SDK
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			"local",
			"local",
			"local",
		)),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	// Create DynamoDB client with local endpoint
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String("http://localhost:8000/")
	})

	// Define table, you can have two primary keys (one for uniquness, one for sorting).
	if !tableExists(context.Background(), client, appConfig.StationTableName) {
		err = createStationsTable(client)
		if err != nil {
			return nil, fmt.Errorf("failed to create station table: %w", err)
		}
	}

	if !tableExists(context.Background(), client, appConfig.TrainTableName) {
		err = createTrainTable(client)
		if err != nil {
			return nil, fmt.Errorf("failed to create train table: %w", err)

		}
	}

	err = insertStations(client)
	if err != nil {
		return nil, fmt.Errorf("failed to insert stations: %w", err)
	}

	return client, nil
}
