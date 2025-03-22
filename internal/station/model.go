package station

import (
	"reflect"
	"strconv"
	"strings"
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
	LineCode1        string  `json:"LineCode1"` // RD, OR, SV, BL, GR
	LineCode2        string  `json:"LineCode2"`
	LineCode3        string  `json:"LineCode3"`
	LineCode4        string  `json:"LineCode4"`
	Longitude        float32 `json:"Lon"`
	Name             string  `json:"Name"`
	StationTogether1 string  `json:"StationTogether1"`
	StationTogether2 string  `json:"StationTogether2"`
}

type GetStationResp struct {
	Address   Address  `json:"address"`
	LineCodes []string `json:"line_codes"`
	Name      string   `json:"station_name"`
}

type DdbStation struct {
	Code             string   `dynamodbav:"code"`
	City             string   `dynamodbav:"city"`
	State            string   `dynamodbav:"state"`
	Street           string   `dynamodbav:"street"`
	Zip              string   `dynamodbav:"zip"`
	Latitude         float32  `dynamodbav:"latitude"`
	Longitude        float32  `dynamodbav:"longitude"`
	Name             string   `dynamodbav:"name"`
	LineCodes        []string `dynamodbav:"lineCodes"`
}

func (s *Station) convertLineCodesToList() []string {
	lineCodes := make([]string, 0, 10)
	for i := 1; i <= 4; i++ {
		lineCodeNum := strconv.Itoa(i)
		r := reflect.ValueOf(s)
		field := reflect.Indirect(r).FieldByName(strings.Join([]string{"LineCode", lineCodeNum}, ""))
		lineCodes = append(lineCodes, field.String())
	}
	return lineCodes
}

func (sl *StationList) ToDdbStations() []DdbStation {
	ddbStations := make([]DdbStation, len(sl.Stations))
	for i, s := range sl.Stations {
		ddbStations[i] = DdbStation{
			State:            s.Address.State,
			City:             s.Address.City,
			Zip:              s.Address.Zip,
			Code:             s.Code,
			Latitude:         s.Latitude,
			LineCodes:        s.convertLineCodesToList(),
			Longitude:        s.Longitude,
			Name:             s.Name,
		}
	}
	return ddbStations
}
