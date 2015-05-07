package delta

import (
	"database/sql"
	//	"bytes"
	//	"encoding/json"
	_ "github.com/mattn/go-oci8"
	//	"io"
	//	"log"
	//	"net/http"
)

type SeismicSite struct {
	Id            int64   `json:"seismic_stream_id"`
	GroundRlnship float64 `json:"ground_rlnship"`
	Height        float64 `json:"height"`
	Latitude      float64 `json:"latitude"`
	Location      string  `json:"location"`
	Longitude     float64 `json:"longitude"`
	Notes         *string `json:"notes"`
	Housing       *string `json:"housing"`
}

func GetSeismicSite(id int64) (*SeismicSite, error) {
	s := SeismicSite{}

	p := "SELECT seismic_site_id, ground_rlnship, height, latitude, location, longitude, notes, housing FROM SEISMIC_SITE WHERE seismic_site_id = :seismic_site_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&s.Id, &s.GroundRlnship, &s.Height, &s.Latitude, &s.Location, &s.Longitude, &s.Notes, &s.Housing)
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

func (s *SeismicSite) GetSeismicStationId() (*int64, error) {
	var id int64

	p := "SELECT seismic_station_id FROM SEISMIC_SITE WHERE seismic_site_id = :seismic_site_id"
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

func (s *SeismicSite) GetSeismicStation() (*SeismicStation, error) {
	id, err := s.GetSeismicStationId()
	if err != nil {
		return nil, err
	} else if id == nil {
		return nil, nil
	}
	return GetSeismicStation(*id)
}
