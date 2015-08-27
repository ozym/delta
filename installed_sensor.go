package delta

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-oci8"
)

type InstalledSensor struct {
	Id             int64     `json:"installed_sensor_id"`
	Azimuth        float64   `json:"azimuth"`
	Dip            float64   `json:"dip"`
	StartTimeStamp time.Time `json:"start_time_stamp"`
	StopTimeStamp  time.Time `json:"stop_time_stamp"`
}

func GetInstalledSensor(id int64) (*InstalledSensor, error) {
	i := InstalledSensor{}

	p := "SELECT installed_sensor_id, azimuth, dip, start_time_stamp, stop_time_stamp FROM installed_sensor WHERE installed_sensor_id = :installed_sensor_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&i.Id, &i.Azimuth, &i.Dip, &i.StartTimeStamp, &i.StopTimeStamp)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &i, nil
}

func GetInstalledSensors() ([]InstalledSensor, error) {
	var installs []InstalledSensor

	q := "SELECT installed_sensor_id, azimuth, dip, start_time_stamp, stop_time_stamp FROM installed_sensor"

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		i := InstalledSensor{}
		if err := rows.Scan(&i.Id, &i.Azimuth, &i.Dip, &i.StartTimeStamp, &i.StopTimeStamp); err != nil {
			return nil, err
		}
		installs = append(installs, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return installs, nil
}

func FindInstalledSensorsBySensorId(id int64) ([]InstalledSensor, error) {
	var installs []InstalledSensor

	p := "SELECT installed_sensor_id, azimuth, dip, start_time_stamp, stop_time_stamp FROM installed_sensor WHERE sensor_id = :sensor_id"
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
		i := InstalledSensor{}
		err = stmt.QueryRow(id).Scan(&i.Id, &i.Azimuth, &i.Dip, &i.StartTimeStamp, &i.StopTimeStamp)
		if err != nil {
			return nil, err
		}
		installs = append(installs, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return installs, nil
}

func (i *InstalledSensor) GetSeismicSiteId() (*int64, error) {
	var id int64

	p := "SELECT seismic_site_id FROM installed_sensor WHERE installed_sensor_id = :installed_sensor_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(i.Id).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &id, nil
}

func (i *InstalledSensor) GetSeismicSite() (*SeismicSite, error) {
	id, err := i.GetSeismicSiteId()
	if err != nil {
		return nil, err
	} else if id == nil {
		return nil, nil
	}
	return GetSeismicSite(*id)
}
