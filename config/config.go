package config

import (
	"encoding/json"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type Configuration struct {
	APIKey       string `json:"api_key"`
	BindingPort  int    `json:"binding_port"`
	IsProd       bool   `json:"prod"`
	TrainRoute   string `json:"train_api"`
	StationRoute string `json:"station_api"`
	APIEndpoint  string `json:"api_endpoint"`
}

func (c *Configuration) GetTrainAPI() string {
	return strings.Join([]string{c.APIEndpoint, c.TrainRoute}, "")
}

func (c *Configuration) GetStationAPI() string {
	return strings.Join([]string{c.APIEndpoint, c.StationRoute}, "")
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
		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&config)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("failed to load config.")
		}
	})
	return &config
}
