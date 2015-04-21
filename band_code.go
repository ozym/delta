package main

import (
	"database/sql"
	//	"bytes"
	//	"encoding/json"
	_ "github.com/mattn/go-oci8"
	//	"io"
	//	"log"
	//	"net/http"
)

type BandCode struct {
	Id           int64   `json:"network_id"`
	Code         string  `json:"code"`
	CornerPeriod *string `json:"corner_period"`
	Description  string  `json:"description"`
	SampleRate   *string `json:"sample_rate"`
}

func GetBandCode(id int64) (*BandCode, error) {
	b := BandCode{}

	p := "SELECT band_code_id, code, corner_period, description, sample_rate FROM BAND_CODE WHERE band_code_id = :band_code_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&b.Id, &b.Code, &b.CornerPeriod, &b.Description, &b.SampleRate)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &b, nil
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
