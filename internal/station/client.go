package station

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/reww406/linetracker/config"
	"github.com/sirupsen/logrus"
)

func requestStations() ([]byte, error) {
	config := config.LoadConfig()
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", config.GetStationAPI(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get request: %w", err)
	}

	req.Header.Add("api_key", config.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get station request: %w", err)
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.WithError(cerr)
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

func GetStations() (*StationList, error) {
	body, err := requestStations()
	if err != nil {
		return nil, fmt.Errorf("requestStations failed with: %w", err)
	}

	log.WithFields(logrus.Fields{
		"Body": string(body),
	}).Info("Got Response from Stations API.")

	var stationList StationList
	if err := json.Unmarshal(body, &stationList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stationList: %w", err)
	}

	return &stationList, nil
}
