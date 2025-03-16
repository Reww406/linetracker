package station

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/reww406/linetracker/config"
)

func GetStations() (*StationList, error) {
  config := config.LoadConfig()
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("api_key", config.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get stations with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var stationList StationList
	if err := json.Unmarshal(body, &stationList); err != nil {
		return nil, err
	}

	return &stationList, nil
}

