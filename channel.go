package delta

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

type Channel struct {
	Id    int64 `json:"channel_id"`
	PinNo int64 `json:"pin_no"`
}

func GetChannel(id int64) (*Channel, error) {
	c := Channel{}

	p := "SELECT channel_id, pin_no FROM CHANNEL WHERE channel_id = :channel_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&c.Id, &c.PinNo)
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

func FindChannelsBySensorId(id int64) ([]Channel, error) {
	var channels []Channel

	p := "SELECT channel_id, pin_no FROM COMPONENT WHERE sensor_id = :sensor_id"

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
		c := Channel{}
		if err := rows.Scan(&c.Id, &c.PinNo); err != nil {
			return nil, err
		}
		channels = append(channels, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return channels, nil
}

func (c *Channel) FindSeismicStreams(start time.Time, stop time.Time) ([]SeismicStream, error) {
	return FindSeismicStreamsByChannelId(c.Id, start, stop)
}
