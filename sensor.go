package delta

import (
	"database/sql"
	//	"bytes"
	//	"encoding/json"
	_ "github.com/mattn/go-oci8"
	//	"io"
	//	"log"
	//	"net/http"
	//	"time"
)

type Sensor struct {
	Id int64 `json:"sensor_id"`
}

func GetSensor(id int64) (*Sensor, error) {
	s := Sensor{}

	p := "SELECT sensor_id FROM SENSOR WHERE sensor_id = :sensor_id"
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

func FindSensorByEquipmentId(id int64) (*Sensor, error) {
	s := Sensor{}

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

func (s *Sensor) FindInstalledSensors() ([]InstalledSensor, error) {
	return FindInstalledSensorsBySensorId(s.Id)
}

func (s *Sensor) FindComponents() ([]Component, error) {
	return FindComponentsBySensorId(s.Id)
}
