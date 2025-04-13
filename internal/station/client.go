package station

import (
	"encoding/json"
	"fmt"

	appConfig "github.com/reww406/linetracker/config"
	"github.com/reww406/linetracker/internal/metro"
	"github.com/sirupsen/logrus"
)

func requestStationTiming(stationCode string) ([]byte, error) {
	config := appConfig.LoadConfig()
	req, err := metro.GetRequest(config.GetStationTimingAPI(stationCode), config.APIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to build request %w", err)
	}
	
	return metro.ExecuteRequest(req)
}

func requestStations() ([]byte, error) {
	config := appConfig.LoadConfig()
	req, err := metro.GetRequest(config.GetStationAPI(), config.APIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to build request %w", err)
	}

	return metro.ExecuteRequest(req)
}

func getStations() (*stationList, error) {
	body, err := requestStations()
	if err != nil {
		return nil, fmt.Errorf("requestStations failed with: %w", err)
	}

	var stations stationList
	if err := json.Unmarshal(body, &stations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stationList: %w", err)
	}

	log.WithFields(logrus.Fields{
		"stations_len": len(stations.Stations),
	}).Info("got response from stations API.")

	return &stations, nil
}


func getStationTimes(stationCode string) (*stationTimeList, error){
	body, err := requestStationTiming(stationCode)
	if err != nil {
		return nil, fmt.Errorf("requestStationTiming failed with %w", err)
	}

	var stationTimes stationTimeList 
	if err := json.Unmarshal(body, &stationTimes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stationList: %w", err)
	}

	log.Info("got stationTimes.")

	return &stationTimes, nil
}
