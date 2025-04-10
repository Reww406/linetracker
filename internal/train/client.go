package train

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/reww406/linetracker/config"
	"github.com/sirupsen/logrus"
)

var log = config.GetLogger()

func requestTrains() ([]byte, error) {
	config := config.LoadConfig()
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", config.GetTrainAPI(), nil)
	if err != nil {
		return nil, fmt.Errorf("get train request could not be created: %w", err)
	}

	req.Header.Add("api_key", config.APIKey)

	logrus.WithFields(logrus.Fields{
		"method":        "GET",
		"train_request": req,
		"api":           config.GetTrainAPI(),
	}).Info("Get Train Request.")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to train API failed: %w", err)
	}

	defer func() {
		cerr := resp.Body.Close()
		if cerr != nil {
			log.WithError(cerr).Error("failed to close train request body.")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get stations with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	return body, nil
}

func GetTrains() (*TrainList, error) {
	body, err := requestTrains()
	if err != nil {
		return nil, fmt.Errorf("failed to get trains from Metro API: %w", err)
	}

	var trainList TrainList
	if err := json.Unmarshal(body, &trainList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal train list: %w", err)
	}

	log.WithField("trains", len(trainList.TrainPredictions)).Info("train predictions returned from API.")

	return &trainList, nil
}
