package train

import (
	"encoding/json"
	"fmt"

	"github.com/reww406/linetracker/config"
	"github.com/reww406/linetracker/internal/metro"
)

var log = config.GetLogger()

func requestTrains() ([]byte, error) {
	config := config.LoadConfig()

	req, err := metro.GetRequest(config.GetTrainAPI(), config.APIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create train req: %w", err)
	}

	return metro.ExecuteRequest(req)
}

func getTrains() (*trainList, error) {
	body, err := requestTrains()
	if err != nil {
		return nil, fmt.Errorf("failed to get trains from Metro API: %w", err)
	}

	var trains trainList
	if err := json.Unmarshal(body, &trains); err != nil {
		return nil, fmt.Errorf("failed to unmarshal train list: %w", err)
	}

	log.WithField("trains", len(trains.TrainPredictions)).Info(
		"train predictions returned from API.",
	)

	return &trains, nil
}
