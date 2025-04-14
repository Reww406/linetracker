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

func itemToDdbStation(item map[string]types.AttributeValue) StationModel {
	longitudeStr := item["longitude"].(*types.AttributeValueMemberN).Value
	longitude, _ := strconv.ParseFloat(longitudeStr, 32)

	latitudeStr := item["latitudeStr"].(*types.AttributeValueMemberN).Value
	latitude, _ := strconv.ParseFloat(latitudeStr, 32)

	lineCodesStr := item["lineCodes"].(*types.AttributeValueMemberSS).Value
	lineCodes := metro.ToLineCodes(lineCodesStr)
	return StationModel{
		Code:      item["code"].(*types.AttributeValueMemberS).Value,
		Name:      item["name"].(*types.AttributeValueMemberS).Value,
		City:      item["city"].(*types.AttributeValueMemberS).Value,
		Zip:       item["zip"].(*types.AttributeValueMemberN).Value,
		Longitude: float32(longitude),
		Latitude:  float32(latitude),
		Street:    item["street"].(*types.AttributeValueMemberS).Value,
		State:     item["state"].(*types.AttributeValueMemberS).Value,
		LineCodes: lineCodes,
	}
}

func ListStations(ctx context.Context, client *dynamodb.Client) (*ListStationResp, error) {
	paginator := dynamodb.NewScanPaginator(client, &dynamodb.ScanInput{
		TableName: config.StationTableName,
	})

	var stations []StationModel
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
			Address: address{
				City:   station.City,
				State:  station.State,
				Street: station.Street,
				Zip:    station.Zip,
			},
			StationCode: station.Code,
			LineCodes:   station.LineCodes,
			Name:        station.Name,
			Schedule:    station.StationSchedule,
		})
	}
	return &ListStationResp{Stations: result}, nil
}
