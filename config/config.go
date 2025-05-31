package config

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Internal configuration for JSON
type jsonConfig struct {
	APIKey             string `json:"api_key"`
	BindingPort        int    `json:"binding_port"`
	IsProd             bool   `json:"prod"`
	TrainRoute         string `json:"train_route"`
	StationRoute       string `json:"station_route"`
	StationTimingRoute string `json:"station_timing_route"`
	APIEndpoint        string `json:"api_endpoint"`
}

type Configuration struct {
	APIKey             string
	BindingPort        int
	IsProd             bool
	trainRoute         string
	stationRoute       string
	stationTimingRoute string
	APIEndpoint        string
	Client             *http.Client
}

func (c *Configuration) GetTrainAPI() string {
	return strings.Join([]string{c.APIEndpoint, c.trainRoute}, "")
}

func (c *Configuration) GetStationAPI() string {
	return strings.Join([]string{c.APIEndpoint, c.stationRoute}, "")
}

func (c *Configuration) GetStationTimingAPI(stationCode string) string {
	return strings.Join([]string{
		c.APIEndpoint, c.stationTimingRoute, stationCode,
	}, "")
}

var (
	config        Configuration
	configureOnce sync.Once
)

func LoadConfig() *Configuration {
	log := GetLogger()
	configureOnce.Do(func() {
		file, err := os.Open("config.json")
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("failed to load config.")
		}
		defer func() {
			cerr := file.Close()
			if cerr != nil {
				log.WithError(cerr).Error("failed to close config file stream.")
			}
		}()

		var j jsonConfig
		decoder := json.NewDecoder(file)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&j)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("failed to load config.")
		}

		config = Configuration{
			APIKey:             j.APIKey,
			BindingPort:        j.BindingPort,
			IsProd:             j.IsProd,
			trainRoute:         j.TrainRoute,
			stationRoute:       j.StationRoute,
			stationTimingRoute: j.StationTimingRoute,
			APIEndpoint:        j.APIEndpoint,
			Client: &http.Client{
				Timeout: 10 * time.Second,
			},
		}
	})

	return &config
}
