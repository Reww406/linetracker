package station

import (
	"database/sql"
	"fmt"
)

var apiURL = "https://api.wmata.com/Rail.svc/json/jStations"

var insertStmt = `
        INSERT OR REPLACE INTO stations (
            code, name, latitude, longitude, line_code1, line_code2, 
            line_code3, line_code4, station_together1, station_together2,
            city, state, street, zip
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

var selectStationByLineCode = `
        SELECT 
            name, line_code1, line_code2, line_code3,
            line_code4, city, state, street, zip
        FROM stations 
        WHERE line_code1 = ? 
            OR line_code2 = ? 
            OR line_code3 = ? 
            OR line_code4 = ?`

func InsertStations(db *sql.DB, stations StationList) error {
	stmt, err := db.Prepare(insertStmt)
	if err != nil {
		return fmt.Errorf("error preparing insert statement: %w", err)
	}
	defer stmt.Close()
	// Insert each station
	for _, station := range stations.Stations {
		_, err := stmt.Exec(
			station.Code,
			station.Name,
			station.Latitude,
			station.Longitude,
			station.LineCode1,
			station.LineCode2,
			station.LineCode3,
			station.LineCode4,
			station.StationTogether1,
			station.StationTogether2,
			station.Address.City,
			station.Address.State,
			station.Address.Street,
			station.Address.Zip,
		)
		if err != nil {
			return fmt.Errorf("error inserting station %s: %w", station.Code, err)
		}
	}
	return nil
}

func createGetStationResp(rows *sql.Rows) (*GetStationResp, error) {
	var (
		name      string
		city      string
		state     string
		street    string
		zip       string
		lineCode1 string
		lineCode2 string
		lineCode3 string
		lineCode4 string
	)
	err := rows.Scan(
		&name,
		&city,
		&state,
		&street,
		&zip,
		&lineCode1,
		&lineCode2,
		&lineCode3,
		&lineCode4,
	)
	if err != nil {
		return nil, err
	}

	address := Address{
		City:   city,
		State:  state,
		Street: street,
		Zip:    zip,
	}
	result := &GetStationResp{
		Name:    name,
		Address: address,

	}
  return result, nil
}

//fmt.Errorf("error scanning station row: %w", err)
func GetStationByLineCode(db *sql.DB, lineCode string) ([]GetStationResp, error) {
	rows, err := db.Query(selectStationByLineCode,
		lineCode, lineCode, lineCode, lineCode)
	if err != nil {
		return nil, fmt.Errorf("error querying stations: %w", err)
	}
	defer rows.Close()
	var result []GetStationResp
	for rows.Next() {
    resp, respErr := createGetStationResp(rows)
		if respErr != nil {
			return nil, fmt.Errorf("error scanning station row: %w", err)
		}
		result = append(result, *resp)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating station rows: %w", err)
	}

	return result, nil
}
