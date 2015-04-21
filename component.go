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

type Component struct {
	Id int64 `json:"component_id"`
}

func GetComponent(id int64) (*Component, error) {
	c := Component{}

	p := "SELECT component_id FROM COMPONENT WHERE component_id = :component_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&c.Id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &c, nil
}

/*
func FindComponentByEquipmentId(id int64) (*Component, error) {
	s := Component{}

	p := "SELECT sensor_id FROM SENSOR WHERE equipment_id = :equipment_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&s.Id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *Component) FindInstalledComponents() ([]InstalledComponent, error) {
	return FindInstalledComponentsByComponentId(s.Id)
}

*/

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

func (c *Component) FindSeismicStreams(start time.Time, stop time.Time) ([]SeismicStream, error) {
	return FindSeismicStreamsByComponentId(c.Id, start, stop)
}

func (c *Component) GetCompOrientId() (*int64, error) {
	var id int64

	p := "SELECT comp_orient_id FROM COMPONENT WHERE component_id = :component_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(c.Id).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &id, nil
}

func (c *Component) GetCompOrient() (*CompOrient, error) {
	id, err := c.GetCompOrientId()
	if err != nil {
		return nil, err
	} else if id == nil {
		return nil, nil
	}
	return GetCompOrient(*id)
}
