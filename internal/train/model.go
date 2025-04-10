package train

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Train struct {
	Car             string `json:"Car"`
	Destination     string `json:"Destination"`
	DestinationCode string `json:"DestinationCode"`
	DestinationName string `json:"DestinationName"`
	Group           string `json:"Group"`
	Line            string `json:"Line"`
	LocationCode    string `json:"LocationCode"`
	LocationName    string `json:"LocationName"`
	Min             string `json:"Min"`
}

type TrainList struct {
	TrainPredictions []Train `json:"Trains"`
}

type DdbTrain struct {
	CarCount        int8   `dynamodbav:"carCount"`
	Destination     string `dynamodbav:"destination"`
	DestinationCode string `dynamodbav:"destinationCode"`
	DestinationName string `dynamodbav:"destinationName"`
	Group           string `dynamodbav:"group"`
	LineCode        string `dynamodbav:"lineCode"`
	LocationCode    string `dynamodbav:"locationCode"`
	LocationName    string `dynamodbav:"locationName"`
	Minutes         int8   `dynamodbav:"minutes"`
	CreatedEpoch    int64  `dynamodbav:"createdEpoch"`
	Id              string `dynamodbav:"id"`
}

func (tl *TrainList) toDdbTrains() []DdbTrain {
	result := make([]DdbTrain, len(tl.TrainPredictions))
	for i, train := range tl.TrainPredictions {
		carInt, _ := strconv.Atoi(train.Car)
		minInt, _ := strconv.Atoi(train.Min)
		result[i] = DdbTrain{
			CarCount:        int8(carInt),
			Destination:     train.Destination,
			DestinationCode: train.DestinationCode,
			DestinationName: train.DestinationName,
			Group:           train.Group,
			LineCode:        train.Line,
			LocationCode:    train.LocationCode,
			LocationName:    train.LocationName,
			Minutes:         int8(minInt),
			CreatedEpoch:    time.Now().Unix(),
			Id:              uuid.New().String(),
		}
	}
	return result
}
