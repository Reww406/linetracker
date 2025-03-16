package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/reww406/linetracker/config"
)

type StationList struct {
	Stations []Station `json:"Stations"`
}

type Address struct {
	City   string `json:"City"`
	State  string `json:"State"`
	Street string `json:"Street"`
	Zip    string `json:"Zip"`
}

type Station struct {
	Address          Address `json:"Address"`
	Code             string  `json:"Code"`
	Latitude         float32 `json:"Lat"`
	LineCode1        string  `json:"LineCode1"`
	LineCode2        string  `json:"LineCode2"`
	LineCode3        string  `json:"LineCode3"`
	LineCode4        string  `json:"LineCode4"`
	Longitude        float32 `json:"Lon"`
	Name             string  `json:"Name"`
	StationTogether1 string  `json:"StationTogether1"`
	StationTogether2 string  `json:"StationTogether2"`
}

var apiURL = "https://api.wmata.com/Rail.svc/json/jStations"

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
