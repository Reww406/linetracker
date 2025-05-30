package train

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"	
)


// TODO Needs to fetch trains every 5 seconds between 6AM->6PM.
// We should implement a Retry
func PollTrainPredictions(client *dynamodb.Client, openClose [2]int64) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Info("Fetching train predictions from Metro API.")

		trainList, err := getTrains()
		if err != nil {
			log.WithError(err).Errorln("failed to get stations from Metro API")
		}

		err = InsertTrains(context.Background(), client, *trainList)
		if err != nil {
			log.WithError(err).Errorln("failed to insert Trains into DDB")
		}
	}
}
