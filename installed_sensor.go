package main

import (
	"database/sql"
	//	"fmt"
	//	"bytes"
	//	"encoding/json"
	_ "github.com/mattn/go-oci8"
	//	"io"
	//	"log"
	//	"net/http"
	"time"
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
	} else if id != nil {
		return nil, nil
	}
	return GetSeismicSite(*id)
}

/*
func GetEquipmentByModel(model string) ([]Equipment, error) {
	var equipment []Equipment

	p := "SELECT e.equipment_id, e.asset_number, e.serial_number, e.notes FROM EQUIPMENT e, EQUIPMENT_MODEL m WHERE e.equipment_model_id = m.equipment_model_id AND m.model = :model"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(model)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := Equipment{}
		if err := rows.Scan(&e.Id, &e.AssetNumber, &e.SerialNumber, &e.Notes); err != nil {
			return nil, err
		}
		equipment = append(equipment, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return equipment, nil
}
*/

/*
func FindEquipmentByModel(model string) ([]Equipment, error) {
	return getEquipments(&model)
}

func FindEquipmentByModels(models []string) ([]Equipment, error) {
	return getEquipments(&models)
}
*/

/*
	p := "SELECT asset_number, serial_number, notes FROM equipment WHERE equipment_id = :equipment_id"
	stmt, err := db.Prepare(p)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&e.AssetNumber, &e.SerialNumber, &e.Notes)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &e, nil
*/

/*
type CsdPair struct {
	ChannelNo int32 `json:"channel_no"`
	PinNo     int32 `json:"pin_no"`
}

func CsdChannels(sensor, pin int32, start time.Time) ([]int32, error) {
	var channels []int32

	sql := fmt.Sprintf("SELECT c.pin_no FROM COMPONENT c, SEISMIC_STREAM s WHERE SEISMIC_STREAM s.EQUIPMENT_MODEL_ID = %d", model)
	fmt.Println(sql)

}

func CsdPairs(model, sensor int32, start time.Time) ([]CsdPair, error) {
	var pairs []CsdPair

	sql := fmt.Sprintf("SELECT o.pin_no FROM COMP_ORIENT o WHERE o.EQUIPMENT_MODEL_ID = %d", model)
	fmt.Println(sql)

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pin int32
		if err := rows.Scan(&pin); err != nil {
			return nil, err
		}
		channels, err := CsdChannels(sensor, pin, start)
		if err != nil {
			return nil, err
		}
		for _, c := range channels {
			pairs = append(pairs, CsdPair{ChannelNo: c, PinNo: pin})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pairs, nil
}

type CsdInstall struct {
	Start    time.Time `json:"start"`
	Channels []CsdPair `json:"channel"`
}

func CsdInstalls(model, sensor int32) ([]CsdInstall, error) {
	var installs []CsdInstall

	sql := fmt.Sprintf("SELECT i.start_time_stamp FROM INSTALLED_SENSOR i WHERE i.SENSOR_ID = %d", sensor)
	fmt.Println(sql)

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var start time.Time
		if err := rows.Scan(&start); err != nil {
			return nil, err
		}
		channels, err := CsdPairs(model, sensor, start)
		if err != nil {
			return nil, err
		}
		if len(channels) > 0 {
			installs = append(installs, CsdInstall{Start: start, Channels: channels})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return installs, nil
}

type Csd struct {
	Model    string       `json:"model"`
	Serial   string       `json:"serial"`
	Installs []CsdInstall `json:"install"`
}

func CsdConfig(indent bool, models []string) ([]byte, error) {
	csds, err := Csds(models)
	if err != nil {
		return nil, err
	}

	if indent {
		return json.MarshalIndent(csds, "", "  ")
	} else {
		return json.Marshal(csds)
	}
}

func Csds(models []string) ([]Csd, error) {
	var csds []Csd

	sql := "SELECT s.sensor_id, m.equipment_model_id, m.model, e.serial_number FROM SENSOR s, EQUIPMENT e, EQUIPMENT_MODEL m WHERE s.EQUIPMENT_ID = e.EQUIPMENT_ID AND e.EQUIPMENT_MODEL_ID = m.EQUIPMENT_MODEL_ID AND m.MODEL IN ("
	for n, m := range models {
		if n > 0 {
			sql += ",'" + m + "'"
		} else {
			sql += "'" + m + "'"
		}
	}
	sql += ")"
	fmt.Println(sql)

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var model, sensor int32
		var name, serial string
		if err := rows.Scan(&sensor, &model, &name, &serial); err != nil {
			return nil, err
		}
		installs, err := CsdInstalls(model, sensor)
		if err != nil {
			return nil, err
		}
		if len(installs) > 0 {
			csds = append(csds, Csd{Model: name, Serial: serial, Installs: installs})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return csds, nil
}
*/

/*
func stations(w http.ResponseWriter, r *http.Request) {
	//var result string = ""

	rows, err := db.Query("SELECT DISTINCT " +
		"       NETWORK.code, " +
		//                "       NETWORK.external_code, " +
		"       SEISMIC_STATION.station_id, " +
		"       SEISMIC_STATION.long_name, " +
		"       PLACE.name, " +
		"       PLACE.latitude, " +
		"       PLACE.longitude, " +
		//                "       EQUIPMENT.asset_number, " +
		"       EQUIPMENT.serial_number, " +
		"       COMPANY.name, " +
		"       EQUIPMENT_MODEL.model " +
		//                "       EQUIPMENT_MODEL.model_nmbr " +
		"FROM " +
		"       network                 NETWORK, " +
		"       sensor                  SENSOR, " +
		"       seismic_station         SEISMIC_STATION, " +
		"       seismic_site            SEISMIC_SITE, " +
		"       installed_sensor        INSTALLED_SENSOR, " +
		"       component               COMPONENT, " +
		"       seismic_stream          SEISMIC_STREAM, " +
		"       channel                 CHANNEL, " +
		"       data_logger             DATA_LOGGER, " +
		"       place                   PLACE, " +
		"       equipment_location      EQUIPMENT_LOCATION, " +
		"       equipment               EQUIPMENT, " +
		"       equipment               EQUIPMENT_PARENT, " +
		"       equipment               SENSOR_EQUIPMENT, " +
		"       equipment_model         EQUIPMENT_MODEL, " +
		"       company                 COMPANY " +
		"WHERE " +
		"       NETWORK.network_id = SEISMIC_STATION.network_id AND " +
		"       SENSOR.sensor_id = INSTALLED_SENSOR.sensor_id AND " +
		"       SEISMIC_STATION.seismic_station_id = SEISMIC_SITE.seismic_station_id AND " +
		"       SEISMIC_SITE.place_id = PLACE.place_id AND " +
		"       INSTALLED_SENSOR.seismic_site_id = SEISMIC_SITE.seismic_site_id AND " +
		"       COMPONENT.sensor_id = INSTALLED_SENSOR.sensor_id AND " +
		"       SEISMIC_STREAM.component_id = COMPONENT.component_id AND " +
		"       SEISMIC_STREAM.channel_id = CHANNEL.channel_id AND " +
		"       CHANNEL.data_logger_id = DATA_LOGGER.data_logger_id AND " +
		"       SENSOR_EQUIPMENT.equipment_id = SENSOR.equipment_id AND " +
		"       EQUIPMENT.equipment_id = DATA_LOGGER.equipment_id AND " +
		"       PLACE.place_id = EQUIPMENT_LOCATION.place_id AND " +
		"       EQUIPMENT_PARENT.equipment_id = EQUIPMENT.equipment_parent_id AND  " +
		"       EQUIPMENT_LOCATION.equipment_id = EQUIPMENT_PARENT.equipment_id AND  " +
		"       EQUIPMENT.equipment_model_id = EQUIPMENT_MODEL.equipment_model_id AND " +
		"       EQUIPMENT_MODEL.company_id = COMPANY.company_id AND " +
		"       EQUIPMENT_LOCATION.installation_date < CURRENT_TIMESTAMP AND " +
		"       EQUIPMENT_LOCATION.removal_date > CURRENT_TIMESTAMP AND " +
		"       SEISMIC_STREAM.start_time_stamp < CURRENT_TIMESTAMP AND " +
		"       SEISMIC_STREAM.stop_time_stamp > CURRENT_TIMESTAMP AND " +
		"       INSTALLED_SENSOR.start_time_stamp < CURRENT_TIMESTAMP AND " +
		"       INSTALLED_SENSOR.stop_time_stamp > CURRENT_TIMESTAMP AND " +
		//                "       EQUIPMENT_MODEL.model_nmbr is not NULL AND " +

		"       SEISMIC_STATION.station_id IN ('TRAB', 'TROB') AND " +
		"       1 = 1 " +
		"ORDER BY " +
		"       SEISMIC_STATION.station_id" +
		"")
*/
/*
   rows, err := db.Query("SELECT DISTINCT " +
           "       NETWORK.code, " +
           //"       NETWORK.external_code, " +
           "       SEISMIC_STATION.station_id, " +
           "       SEISMIC_STATION.long_name, " +
           "       PLACE.name, " +
           "       PLACE.latitude, " +
           "       PLACE.longitude, " +
           //"       EQUIPMENT.asset_number, " +
           "       EQUIPMENT.serial_number, " +
           "       COMPANY.name, " +
           "       EQUIPMENT_MODEL.model, " +
           //"       EQUIPMENT_MODEL.model_nmbr, " +
           "       DATA_LOGGER.short_sn " +
           "FROM " +
           "       network                 NETWORK, " +
           "       sensor                  SENSOR, " +
           "       seismic_station         SEISMIC_STATION, " +
           "       seismic_site            SEISMIC_SITE, " +
           "       installed_sensor        INSTALLED_SENSOR, " +
           "       component               COMPONENT, " +
           "       seismic_stream          SEISMIC_STREAM, " +
           "       channel                 CHANNEL, " +
           "       data_logger             DATA_LOGGER, " +
           "       place                   PLACE, " +
           "       equipment_location      EQUIPMENT_LOCATION, " +
           "       equipment               EQUIPMENT, " +
           "       equipment_model         EQUIPMENT_MODEL, " +
           "       company                 COMPANY " +

           "WHERE " +
           "       NETWORK.network_id = SEISMIC_STATION.network_id AND " +
           "       SENSOR.sensor_id = INSTALLED_SENSOR.sensor_id AND " +
           "       SEISMIC_STATION.seismic_station_id = SEISMIC_SITE.seismic_station_id AND " +
           "       SEISMIC_SITE.place_id = PLACE.place_id AND " +
           "       INSTALLED_SENSOR.seismic_site_id = SEISMIC_SITE.seismic_site_id AND " +
           "       COMPONENT.sensor_id = INSTALLED_SENSOR.sensor_id AND " +
           "       SEISMIC_STREAM.component_id = COMPONENT.component_id AND " +
           "       SEISMIC_STREAM.channel_id = CHANNEL.channel_id AND " +
           "       CHANNEL.data_logger_id = DATA_LOGGER.data_logger_id AND " +
           "       EQUIPMENT.equipment_id = DATA_LOGGER.equipment_id AND " +
           "       PLACE.place_id = EQUIPMENT_LOCATION.place_id AND " +
           "       EQUIPMENT_LOCATION.equipment_id = EQUIPMENT.equipment_id AND " +
           "       EQUIPMENT.equipment_model_id = EQUIPMENT_MODEL.equipment_model_id AND " +
           "       EQUIPMENT_MODEL.company_id = COMPANY.company_id AND " +
           "       EQUIPMENT_LOCATION.installation_date < CURRENT_TIMESTAMP AND " +
           "       EQUIPMENT_LOCATION.removal_date > CURRENT_TIMESTAMP AND " +
           "       SEISMIC_STREAM.start_time_stamp < CURRENT_TIMESTAMP AND " +
           "       SEISMIC_STREAM.stop_time_stamp > CURRENT_TIMESTAMP AND " +
           "       INSTALLED_SENSOR.start_time_stamp < CURRENT_TIMESTAMP AND " +
           "       INSTALLED_SENSOR.stop_time_stamp > CURRENT_TIMESTAMP AND " +
           //"       EQUIPMENT_MODEL.model_nmbr is not NULL AND " +
           "       DATA_LOGGER.short_sn is not NULL AND " +
           "       1 = 1 " +
           "ORDER BY " +
           "       SEISMIC_STATION.station_id" +
           "");
*/
/*

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var stations map[string]Station

	stations = make(map[string]Station)
	for rows.Next() {
		var network string
		//var external string
		var station string
		var name string
		var place string
		var lat float32
		var lon float32
		//var asset string
		var serial string
		var company string
		var model string
		//var number string
		//var logger string

		if err := rows.Scan(&network, &station, &name, &place, &lat, &lon, &serial, &company, &model); err != nil {
			log.Fatal(err)
		}
		stations[name] = Station{network, station, name, place, lat, lon, serial, company, model}
		//result = result + fmt.Sprintf(" %s %g %g", name, lat, lon)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(stations)
	if err != nil {
		log.Fatal(err)
	}

	io.WriteString(w, bytes.NewBuffer(b).String())
}
*/
