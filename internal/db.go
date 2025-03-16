package internal

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/reww406/linetracker/config"
	"github.com/sirupsen/logrus"
)

var log = config.GetLogger()

func stationsExists(db *sql.DB) bool {
	stmt, _ := db.Prepare("SELECT EXISTS(SELECT 1 FROM stations LIMIT 1)")
	defer stmt.Close()

	var exists bool
	err := stmt.QueryRow().Scan(&exists)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to search for if stations exists")
		return false
	}
	return exists
}

func insertStations(db *sql.DB, stations StationList) error {
	stmt, err := db.Prepare(`
        INSERT OR REPLACE INTO stations (
            code, name, latitude, longitude, line_code1, line_code2, 
            line_code3, line_code4, station_together1, station_together2,
            city, state, street, zip
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `)
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

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	createStationTableSQL := `
	CREATE TABLE IF NOT EXISTS stations (
		code TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		latitude REAL NOT NULL,
		longitude REAL NOT NULL,
		line_code1 TEXT,
		line_code2 TEXT,
		line_code3 TEXT,
		line_code4 TEXT,
		station_together1 TEXT NOT NULL,
		station_together2 TEXT NOT NULL,
		city TEXT NOT NULL,
		state TEXT NOT NULL,
		street TEXT NOT NULL,
		zip TEXT NOT NULL
	);`

	_, err = db.Exec(createStationTableSQL)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error creating stations table: %w", err)
	}

	if !stationsExists(db) {
		stations, err := GetStations()
		if err != nil {
		  return nil, err	
		}
		if insertStations(db, *stations) != nil {
      return nil , err
		}
	}
	return db, nil
}
