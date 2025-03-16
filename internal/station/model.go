package station

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
