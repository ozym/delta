package main

import (
	"database/sql"
	//	"bytes"
	//	"encoding/json"
	_ "github.com/mattn/go-oci8"
	//	"io"
	//	"log"
	//	"net/http"
	"time"
)

type SeismicStation struct {
	Id            int64     `json:"seismic_station_id"`
	DateClosed    time.Time `json:"date_closed"`
	DateOpened    time.Time `json:"date_opened"`
	FileReference *string   `json:"file_reference"`
	Height        float64   `json:"height"`
	Latitude      float64   `json:"latitude"`
	LongName      *string   `json:"long_name"`
	Longitude     float64   `json:"longitude"`
	Notes         *string   `json:"notes"`
	StationId     string    `json:"station_id"`
}

func GetSeismicStation(id int64) (*SeismicStation, error) {
	s := SeismicStation{}

	p := "SELECT seismic_station_id, date_closed, date_opened, file_reference, height, latitude, long_name, longitude, notes, station_id FROM SEISMIC_STATION WHERE seismic_station_id = :seismic_station_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&s.Id, &s.DateClosed, &s.DateOpened, &s.FileReference, &s.Height, &s.Latitude, &s.LongName, &s.Longitude, &s.Notes, &s.StationId)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &s, nil
}

/*
func (s *Component) FindInstalledComponents() ([]InstalledComponent, error) {
	return FindInstalledComponentsByComponentId(s.Id)
}

*/

/*
func FindComponentsBySensorId(id int64) ([]Component, error) {
	var components []Component

	p := "SELECT component_id FROM COMPONENT WHERE sensor_id = :sensor_id"

	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		c := Component{}
		if err := rows.Scan(&c.Id); err != nil {
			return nil, err
		}
		components = append(components, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return components, nil
}
*/

/*
func (c *Component) FindSeismicStreams(start time.Time, stop time.Time) ([]SeismicStream, error) {
	return FindSeismicStreamsByComponentId(c.Id, start, stop)
}
*/

func (s *SeismicStation) GetNetworkId() (*int64, error) {
	var id int64

	p := "SELECT network_id FROM SEISMIC_STATION WHERE seismic_station_id = :seismic_station_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(s.Id).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &id, nil
}

func (s *SeismicStation) GetNetwork() (*Network, error) {
	id, err := s.GetNetworkId()
	if err != nil {
		return nil, err
	} else if id == nil {
		return nil, nil
	}
	return GetNetwork(*id)
}
