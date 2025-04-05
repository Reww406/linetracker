package station

import (
	"encoding/json"
	"testing"
)

func TestStationList(t *testing.T) {
	stationList := StationList{
		Stations: []Station{
			{
				Address: Address{
					City:   "Test",
					State:  "Test",
					Street: "Test",
					Zip:    "22030",
				},
				Code:             "Test",
				Latitude:         32.3,
				Longitude:        43.2,
				LineCode1:        "Test1",
				LineCode2:        "Test2",
				LineCode3:        "Test3",
				LineCode4:        "Test4",
				StationTogether1: "Test",
				StationTogether2: "Test",
			},
		},
	}

	stationList.ToDdbStations()
}

func TestMarshallingToStationList(t *testing.T) {
	input := `{
		"Stations": [{
			"Code": "A01",
			"Name": "Metro Center",
			"StationTogether1": "C01",
			"StationTogether2": "",
			"LineCode1": "RD",
			"LineCode2": null,
			"LineCode3": null,
			"LineCode4": null,
			"Lat": 38.898303,
			"Lon": -77.028099,
			"Address": {
				"Street": "607 13th St. NW",
				"City": "Washington",
				"State": "DC",
				"Zip": "20005"
			}
		}]
	}`
	var stationList StationList
	if err := json.Unmarshal([]byte(input), &stationList); err != nil {
		log.Fatal(err.Error())	
	}
}
