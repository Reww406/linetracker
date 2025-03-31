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

func tableExists(ctx context.Context, client *dynamodb.Client, tableName string) bool {
	_, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	return err == nil
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
	tableName := "stations"
	if !tableExists(context.Background(), client, tableName) {

		log.WithFields(logrus.Fields{
			"TableName": tableName,
		}).Info("Creating DDB Table")
		_, err = client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
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
			return nil, fmt.Errorf("error creating stations table: %w", err)
		}

		stations, err := station.GetStations()
		if err != nil {
			return nil, err
		}
		if err := station.InsertStations(context.Background(), client, *stations); err != nil {
			return nil, err
		}
	}

	return client, nil
}
