package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-oci8"
)

type Impact struct {
	Name      string  `json:"Name"`
	Rate      float32 `json:"Rate"`
	Gain      float32 `json:"Gain"`
	Q         float32 `json:"Q"`
	Latitude  float32 `json:"Latitude"`
	Longitude float32 `json:"Longitude"`
}

var Q = map[int32]float32{200: 0.98829, 100: 0.97671, 50: 0.95395}

func ImpactConfig(indent bool) ([]byte, error) {
	impacts, err := Impacts()
	if err != nil {
		return nil, err
	}

	if indent {
		j, err := json.MarshalIndent(impacts, "", "  ")
		if err != nil {
			return nil, err
		}
		return append(j, '\n'), nil
	} else {
		return json.Marshal(impacts)
	}

}

func Impacts() (map[string]Impact, error) {
	var impacts map[string]Impact

	impacts = make(map[string]Impact)

	var queries []string

	// gather seismic stream information ...
	queries = append(queries, "SELECT DISTINCT "+
		" SEISMIC_STATION.long_name,"+
		" SEISMIC_STATION.latitude,"+
		" SEISMIC_STATION.longitude,"+
		" SEISMIC_STATION.station_id,"+
		" SEISMIC_SITE.location,"+
		" CONCAT(BAND_CODE.code, CONCAT(SENSOR_CODE.code, ORIENT_CODE.code)),"+
		" NETWORK.code,"+
		" NETWORK.external_code,"+
		" SEISMIC_STREAM.sample_rate,"+
		" SEISMIC_STREAM.sensitivity"+
		" FROM seismic_stream SEISMIC_STREAM"+
		" JOIN band_code BAND_CODE ON BAND_CODE.band_code_id = SEISMIC_STREAM.band_code_id"+
		" JOIN sensor_code SENSOR_CODE ON SENSOR_CODE.sensor_code_id = SEISMIC_STREAM.sensor_code_id"+
		" JOIN orient_code ORIENT_CODE ON ORIENT_CODE.orient_code_id = SEISMIC_STREAM.orient_code_id"+
		" JOIN component COMPONENT ON COMPONENT.component_id = SEISMIC_STREAM.component_id"+
		" JOIN sensor SENSOR ON SENSOR.sensor_id = COMPONENT.sensor_id"+
		" LEFT JOIN installed_sensor INSTALLED_SENSORS ON INSTALLED_SENSORS.sensor_id = SENSOR.sensor_id"+
		" LEFT JOIN seismic_site SEISMIC_SITE ON SEISMIC_SITE.seismic_site_id = INSTALLED_SENSORS.seismic_site_id"+
		" LEFT JOIN seismic_station SEISMIC_STATION ON SEISMIC_STATION.seismic_station_id = SEISMIC_SITE.seismic_station_id"+
		" LEFT JOIN network NETWORK ON NETWORK.network_id = SEISMIC_STATION.network_id"+
		" JOIN equipment equipment_2 ON equipment_2.equipment_id = sensor.equipment_id"+
		" JOIN equipment_model equipment_model_2 ON equipment_model_2.equipment_model_id = equipment_2.equipment_model_id"+
		" WHERE INSTALLED_SENSORS.start_time_stamp <= SEISMIC_STREAM.start_time_stamp"+
		" AND INSTALLED_SENSORS.stop_time_stamp > CURRENT_TIMESTAMP"+
		" AND SEISMIC_STREAM.start_time_stamp < CURRENT_TIMESTAMP"+
		" AND SEISMIC_STREAM.stop_time_stamp > CURRENT_TIMESTAMP"+
		" AND CONCAT(BAND_CODE.code, CONCAT(SENSOR_CODE.code, ORIENT_CODE.code))"+
		" IN ('HHZ', 'HHN', 'HHE', 'HH1', 'HH2', 'EHZ', 'EHN', 'EHE', 'EH1', 'EH2',"+
		" 'HNZ', 'HNN', 'HNE', 'HN1', 'HN2', 'BNZ', 'BNN', 'BNE', 'BN1', 'BN2')"+
		"")

	for _, sql := range queries {
		rows, err := db.Query(sql)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {

			var name *string
			var lat, lon float32
			var sta, loc, cha, net string
			var ext *string
			var rate int32
			var gain float32

			if err := rows.Scan(&name, &lat, &lon, &sta, &loc, &cha, &net, &ext, &rate, &gain); err != nil {
				return nil, err
			}
			if name == nil {
				continue
			}

			q, ok := Q[rate]
			if ok {
				var srcname string

				if ext != nil {
					srcname = fmt.Sprintf("%s_%s_%s_%s", *ext, sta, loc, cha)
				} else {
					srcname = fmt.Sprintf("%s_%s_%s_%s", net, sta, loc, cha)
				}

				impacts[srcname] = Impact{
					Name:      *name,
					Rate:      (float32)(rate),
					Gain:      gain,
					Q:         q,
					Latitude:  lat,
					Longitude: lon,
				}
			}
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return impacts, nil
}
