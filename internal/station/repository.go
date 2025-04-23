package station

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/reww406/linetracker/config"
	"github.com/reww406/linetracker/internal/metro"
	"github.com/sirupsen/logrus"
)

var log = config.GetLogger()

func createDdbStations(stationList stationList) ([]StationModel, error) {
	// will tick 5 times a seconds.
	limiter := time.NewTicker(200 * time.Millisecond)
	defer limiter.Stop()

	result := make([]StationModel, len(stationList.Stations))
	for i, station := range stationList.Stations {
		// wait for limiter to deliever on a channel before running
		<-limiter.C
		stationTimes, err := getStationTimes(station.Code)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get station times for station: %s with error: %w", station.Code, err,
			)
		}
		result[i] = station.toStationModel(*stationTimes)
	}
	return result, nil
}

func InsertStations(ctx context.Context, client *dynamodb.Client) error {
	stationList, err := getStations()
	if err != nil {
		return fmt.Errorf("failed to get stations: %w", err)
	}

	ddbStations, err := createDdbStations(*stationList)
	if err != nil {
		return err
	}

	log.WithFields(logrus.Fields{
		"stations_len": len(ddbStations),
	}).Info("inserting stations into DDB")

	for _, station := range ddbStations {
		item, err := attributevalue.MarshalMap(station)
		if err != nil {
			return fmt.Errorf("failed to marshal station: %w", err)
		}

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: config.StationTableName,
			Item:      item,
		})
		if err != nil {
			return fmt.Errorf("failed to insert station %s: %w", station.Code, err)
		}
	}
	return nil
}

func ddbListToStringList(values []types.AttributeValue) []string {
	result := make([]string, len(values))
	for i, val := range values {
		// TODO What the hell is happening here.
		if sv, ok := val.(*types.AttributeValueMemberS); ok {
			result[i] = sv.Value
		}
	}
	return result
}

func itemToDdbStation(item map[string]types.AttributeValue) StationModel {
	longitudeStr := item["longitude"].(*types.AttributeValueMemberN).Value
	longitude, _ := strconv.ParseFloat(longitudeStr, 32)

	latitudeStr := item["latitude"].(*types.AttributeValueMemberN).Value
	latitude, _ := strconv.ParseFloat(latitudeStr, 32)

	lineCodesStr := item["lineCodes"].(*types.AttributeValueMemberL).Value
	lineCodes := metro.ToLineCodes(lineCodesStr)

	return StationModel{
		Code:         item["code"].(*types.AttributeValueMemberS).Value,
		Name:         item["name"].(*types.AttributeValueMemberS).Value,
		City:         item["city"].(*types.AttributeValueMemberS).Value,
		Zip:          item["zip"].(*types.AttributeValueMemberS).Value,
		Longitude:    float32(longitude),
		Latitude:     float32(latitude),
		Street:       item["street"].(*types.AttributeValueMemberS).Value,
		State:        item["state"].(*types.AttributeValueMemberS).Value,
		LineCodes:    lineCodes,
		Destinations: ddbListToStringList(item["destinations"].(*types.AttributeValueMemberL).Value),
	}
}

func createStationCodeLookup(stations []StationModel) map[string]StationModel {
	result := make(map[string]StationModel, len(stations))
	for _, station := range stations {
		result[station.Code] = station
	}
	return result
}

func scanStations(ctx context.Context, client *dynamodb.Client) ([]StationModel, error) {
	paginator := dynamodb.NewScanPaginator(client, &dynamodb.ScanInput{
		TableName: config.StationTableName,
	})

	var result []StationModel

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failled to scan stations table: %w", err)
		}
		for _, item := range page.Items {
			result = append(result, itemToDdbStation(item))
		}
	}
	log.WithField("stationsFound", len(result)).Info("stations found from DDB.")

	return result, nil
}

func ListStations(ctx context.Context, client *dynamodb.Client) ([]StationModel, error) {
	stations, err := scanStations(ctx, client)
	if err != nil {
		return nil, err
	}
	return stations, nil
}

func GetDestinationStations(ctx context.Context, client *dynamodb.Client) ([]StationModel, error) {
	stations, err := scanStations(ctx, client)
	stationCodeLookup := createStationCodeLookup(stations)
	if err != nil {
		return nil, err
	}
	set := make(map[string]StationModel)
	for _, station := range stations {
		for _, destination := range station.Destinations {
			set[destination] = stationCodeLookup[destination]
		}
	}

	result := make([]StationModel, 0, len(set))
	for _, v := range set {
		result = append(result, v)
	}

	return result, nil
}
