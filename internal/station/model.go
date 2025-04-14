package station

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/reww406/linetracker/internal/metro"
)

type stationList struct {
	Stations []stationData `json:"Stations"`
}

type address struct {
	City   string `json:"City"`
	State  string `json:"State"`
	Street string `json:"Street"`
	Zip    string `json:"Zip"`
}

type stationData struct {
	Address          address        `json:"Address"`
	Code             string         `json:"Code"`
	Latitude         float32        `json:"Lat"`
	LineCode1        metro.LineCode `json:"LineCode1"`
	LineCode2        metro.LineCode `json:"LineCode2"`
	LineCode3        metro.LineCode `json:"LineCode3"`
	LineCode4        metro.LineCode `json:"LineCode4"`
	Longitude        float32        `json:"Lon"`
	Name             string         `json:"Name"`
	StationTogether1 string         `json:"StationTogether1"`
	StationTogether2 string         `json:"StationTogether2"`
}

type stationTimeList struct {
	StationTimes []stationTimesData `json:"StationTimes"`
}

type stationTimesData struct {
	Code        string      `json:"Code"`
	StationName string      `json:"StationName"`
	Monday      daySchedule `json:"Monday"`
	Tuesday     daySchedule `json:"Tuesday"`
	Wednesday   daySchedule `json:"Wednesday"`
	Thursday    daySchedule `json:"Thursday"`
	Friday      daySchedule `json:"Friday"`
	Saturday    daySchedule `json:"Saturday"`
	Sunday      daySchedule `json:"Sunday"`
}

type daySchedule struct {
	OpeningTime string  `json:"OpeningTime"`
	FirstTrains []train `json:"FirstTrains"`
	LastTrains  []train `json:"LastTrains"`
}

type train struct {
	LeavingTime        string `json:"LeavingTime"`
	DestinationStation string `json:"DestinationStation"`
}

type ListStationResp struct {
	Stations []GetStationResp `json:"stations"`
}

type GetStationResp struct {
	Address     address           `json:"address"`
	LineCodes   []metro.LineCode  `json:"line_codes"`
	StationCode string            `json:"station_code"`
	Name        string            `json:"name"`
	Schedule    []StationSchedule `json:"schedule"`
}

type StationModel struct {
	Code            string            `dynamodbav:"code"`
	City            string            `dynamodbav:"city"`
	State           string            `dynamodbav:"state"`
	Street          string            `dynamodbav:"street"`
	Zip             string            `dynamodbav:"zip"`
	Latitude        float32           `dynamodbav:"latitude"`
	Longitude       float32           `dynamodbav:"longitude"`
	Name            string            `dynamodbav:"name"`
	LineCodes       []metro.LineCode        `dynamodbav:"lineCodes"`
	StationSchedule []StationSchedule `dynamodbav:"stationSchedule"`
}

type StationSchedule struct {
	Day         string `dynamodbav:"day"`
	OpeningTime string `dynamodbav:"openingTime"`
	LastTrain   string `dynamodbav:"lastTrain"`
}

func (s *stationData) convertLineCodesToList() []metro.LineCode {
	result := make([]metro.LineCode, 0, 4)
	for i := 1; i <= 4; i++ {
		lineCodeNum := strconv.Itoa(i)
		r := reflect.ValueOf(s)
		field := reflect.Indirect(r).FieldByName(strings.Join([]string{"LineCode", lineCodeNum}, ""))
		if field.String() != "" {
			result = append(result, metro.LineCode(field.String()))
		}
	}
	return result
}

func (st *stationTimeList) toStationSchedule() ([]StationSchedule, error) {
	days := []string{
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
		"Sunday",
	}
	if len(st.StationTimes) != 1 {
		return nil, fmt.Errorf("station times had more than one entry: %d", len(st.StationTimes))
	}

	result := make([]StationSchedule, len(days))
	r := reflect.ValueOf(st.StationTimes[0])
	for i, day := range days {
		field := reflect.Indirect(r).FieldByName(day)
		daySchedule := field.Interface().(daySchedule)

		lastTrain := ""
		if len(daySchedule.LastTrains) != 0 {
			lastTrain = daySchedule.LastTrains[0].LeavingTime
		}

		result[i] = StationSchedule{
			Day:         day,
			OpeningTime: daySchedule.OpeningTime,
			LastTrain:   lastTrain,
		}
	}
	return result, nil
}

func (s *stationData) toStationModel(stationTimes stationTimeList) StationModel {
	daySchedules, err := stationTimes.toStationSchedule()
	if err != nil {
		log.WithError(err).Errorln("failed to covert stationTimes to DdbDaySchedule")
		daySchedules = []StationSchedule{}
	}

	return StationModel{
		State:           s.Address.State,
		City:            s.Address.City,
		Zip:             s.Address.Zip,
		Code:            s.Code,
		Latitude:        s.Latitude,
		LineCodes:       s.convertLineCodesToList(),
		Longitude:       s.Longitude,
		Name:            s.Name,
		StationSchedule: daySchedules,
	}
}
