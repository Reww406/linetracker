package config

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type Configuration struct {
	APIKey      string `json:"api_key"`
	BindingPort int    `json:"binding_port"`
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
			}).Fatal("failed to log config.")
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&config)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("failed to log config.")
		}
	})
	return &config
}
