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

type Network struct {
	Id           int64   `json:"network_id"`
	AlbumPostfix *string `json:"album_postfix"`
	Code         string  `json:"code"`
	Description  string  `json:"description"`
	ExternalCode *string `json:"external_code"`
	PublicFlag   *bool   `json:"public_flag"`
}

func GetNetwork(id int64) (*Network, error) {
	n := Network{}

	p := "SELECT network_id, album_postfix, code, description, external_code, public_flag FROM NETWORK WHERE network_id = :network_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&n.Id, &n.AlbumPostfix, &n.Code, &n.Description, &n.ExternalCode, &n.PublicFlag)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &n, nil
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
