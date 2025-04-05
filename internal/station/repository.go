package station

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/reww406/linetracker/config"
	"github.com/sirupsen/logrus"
)

var log = config.GetLogger()

func itemToDdbStation(item map[string]types.AttributeValue) DdbStation {	
	longitudeStr := item["longitude"].(*types.AttributeValueMemberN).Value
	longitude, _ := strconv.ParseFloat(longitudeStr, 32)

	latitudeStr := item["latitudeStr"].(*types.AttributeValueMemberN).Value
	latitude, _ := strconv.ParseFloat(latitudeStr, 32)
	return DdbStation{
		Code:      item["code"].(*types.AttributeValueMemberS).Value,
		Name:      item["name"].(*types.AttributeValueMemberS).Value,
		City:      item["city"].(*types.AttributeValueMemberS).Value,
		Zip:       item["zip"].(*types.AttributeValueMemberN).Value,
		Longitude: float32(longitude),
		Latitude:  float32(latitude),
		Street:    item["street"].(*types.AttributeValueMemberS).Value,
		State:     item["state"].(*types.AttributeValueMemberS).Value,
		LineCodes: item["lineCodes"].(*types.AttributeValueMemberSS).Value,
	}
}

func InsertStations(ctx context.Context, client *dynamodb.Client, stationList StationList) error {
	ddbStations := stationList.ToDdbStations()
	log.WithFields(logrus.Fields{
		"StationsToInsert": len(ddbStations),
	}).Info("Inserting Stations into DDB")

	for _, station := range ddbStations {
		item, err := attributevalue.MarshalMap(station)
		if err != nil {
			return fmt.Errorf("failed to marshal station: %w", err)
		}

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String("stations"),
			Item:      item,
		})
		if err != nil {
			return fmt.Errorf("failed to insert station %s: %w", station.Code, err)
		}
	}
	return nil
}

func ListStations(ctx context.Context, client *dynamodb.Client) (*ListStationResp, error) {
	paginator := dynamodb.NewScanPaginator(client, &dynamodb.ScanInput{
		TableName: aws.String("stations"),
	})

	var stations []DdbStation
	var result []GetStationResp

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failled to scan stations table: %w", err)
		}
		for _, item := range page.Items {
			stations = append(stations, itemToDdbStation(item))
		}
	}

	for _, station := range stations {
		result = append(result, GetStationResp{
			Address: Address{
				City:   station.City,
				State:  station.State,
				Street: station.Street,
				Zip:    station.Zip,
			},
			StationCode: station.Code,
			LineCodes:   station.LineCodes,
			Name:        station.Name,
		})
	}
	return &ListStationResp{Stations: result}, nil
}
