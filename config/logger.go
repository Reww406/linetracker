package config

import (
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
	once   sync.Once
)

func GetLogger() *logrus.Logger {
	once.Do(func() {
		log = logrus.New()

		file, err := os.OpenFile("/logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			logrus.Fatal(err)
		}

		log.SetOutput(io.MultiWriter(file, os.Stdout))
		log.SetFormatter(&logrus.JSONFormatter{})
		log.SetLevel(logrus.DebugLevel)
	})
	return log
}
