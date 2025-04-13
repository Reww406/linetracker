package station

import (
	"encoding/json"
	"testing"
)

var stationTimes = stationTimeList{
    StationTimes: []stationTimesData{
        {
            Code:        "E10",
            StationName: "Greenbelt",
            Monday: daySchedule{
                OpeningTime: "04:50",
                FirstTrains: []train{
                    {LeavingTime: "05:00", DestinationStation: "F11"},
                },
                LastTrains: []train{
                    {LeavingTime: "23:26", DestinationStation: "F11"},
                },
            },
            Tuesday: daySchedule{
                OpeningTime: "04:50",
                FirstTrains: []train{
                    {LeavingTime: "05:00", DestinationStation: "F11"},
                },
                LastTrains: []train{
                    {LeavingTime: "23:26", DestinationStation: "F11"},
                },
            },
            Wednesday: daySchedule{
                OpeningTime: "04:50",
                FirstTrains: []train{
                    {LeavingTime: "05:00", DestinationStation: "F11"},
                },
                LastTrains: []train{
                    {LeavingTime: "23:26", DestinationStation: "F11"},
                },
            },
            Thursday: daySchedule{
                OpeningTime: "04:50",
                FirstTrains: []train{
                    {LeavingTime: "05:00", DestinationStation: "F11"},
                },
                LastTrains: []train{
                    {LeavingTime: "23:26", DestinationStation: "F11"},
                },
            },
            Friday: daySchedule{
                OpeningTime: "04:50",
                FirstTrains: []train{
                    {LeavingTime: "05:00", DestinationStation: "F11"},
                },
                LastTrains: []train{
                    {LeavingTime: "00:26", DestinationStation: "F11"},
                },
            },
            Saturday: daySchedule{
                OpeningTime: "06:50",
                FirstTrains: []train{
                    {LeavingTime: "07:00", DestinationStation: "F11"},
                },
                LastTrains: []train{
                    {LeavingTime: "00:26", DestinationStation: "F11"},
                },
            },
            Sunday: daySchedule{
                OpeningTime: "07:50",
                FirstTrains: []train{
                    {LeavingTime: "07:00", DestinationStation: "F11"},
                },
                LastTrains: []train{
                    {LeavingTime: "23:26", DestinationStation: "F11"},
                },
            },
        },
    },
}

func TestStationList(t *testing.T) {
	stations := stationList{
		Stations: []stationData{
			{
				Address: address{
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

	stations.Stations[0].toStationModel(stationTimes)
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
	var stations stationList
	if err := json.Unmarshal([]byte(input), &stations); err != nil {
		log.Fatal(err.Error())	
	}
}
