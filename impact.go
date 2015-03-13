package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"io"
	"log"
	"net/http"
	//        "time"
	//        "strings"
)

/*
type Pairs struct {
        Loggers []string `json:"loggers"`
        Sensors []string `json:"sensors"`
}
*/

type Impact struct {
	Name      string
	Rate      float32
	Gain      float32
	Q         float32
	Latitude  float32
	Longitude float32
}

var Q = map[int32]float32{200: 0.98829, 100: 0.97671, 50: 0.95395}

func Impacts() map[string]Impact {
	//var impacts []Impact
	var impacts map[string]Impact

	impacts = make(map[string]Impact)

	var queries []string

	// first check the sensor installation times and location times are consistent ...
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

	/*
		queries = append(queries, "SELECT DISTINCT "+
			"       SEISMIC_STATION.station_id, "+
			"       SEISMIC_STATION.long_name, "+
			"       SEISMIC_SITE.location, "+
			//"       COMP_ORIENT.pin_no, " +
			"       INSTALLED_SENSOR.start_time_stamp, "+
			"       INSTALLED_SENSOR.stop_time_stamp, "+
			"       SENSOR_LOCATION.installation_date, "+
			"       SENSOR_LOCATION.removal_date, "+
			"       SENSOR_MODEL.model, "+
			"       SENSOR_EQUIPMENT.serial_number "+
			"FROM "+
			"       sensor                  SENSOR, "+
			"       seismic_station         SEISMIC_STATION, "+
			"       seismic_site            SEISMIC_SITE, "+
			"       installed_sensor        INSTALLED_SENSOR, "+
			//"       component               COMPONENT, " +
			//"       comp_orient             COMP_ORIENT, " +
			"       equipment_location      SENSOR_LOCATION, "+
			"       equipment               SENSOR_EQUIPMENT, "+
			"       equipment_model         SENSOR_MODEL "+
			"WHERE "+
			"       SENSOR.sensor_id = INSTALLED_SENSOR.sensor_id AND "+
			"       INSTALLED_SENSOR.seismic_site_id = SEISMIC_SITE.seismic_site_id AND "+
			"       SEISMIC_SITE.seismic_station_id = SEISMIC_STATION.seismic_station_id AND "+
			"       SENSOR_EQUIPMENT.equipment_id = SENSOR.equipment_id AND "+
			"       SENSOR_MODEL.equipment_model_id = SENSOR_EQUIPMENT.equipment_model_id AND "+
			"       SENSOR_LOCATION.equipment_id = SENSOR_EQUIPMENT.equipment_id AND  "+
			//"       COMPONENT.sensor_id = SENSOR.sensor_id AND  " +
			//"       COMPONENT.comp_orient_id = COMP_ORIENT.comp_orient_id AND  " +
			//                "       NOT EXISTS(SELECT * FROM equipment e WHERE e.equipment_parent_id = SENSOR_EQUIPMENT.equipment_id) AND " +
			//                "       EQUIPMENT.equipment_id = DATA_LOGGER.equipment_id AND " +
			//                "       PLACE.place_id = EQUIPMENT_LOCATION.place_id AND " +
			//                "       EQUIPMENT_PARENT.equipment_id = EQUIPMENT.equipment_parent_id AND  " +
			//                "       EQUIPMENT_LOCATION.equipment_id = EQUIPMENT_PARENT.equipment_id AND  " +
			//                "       EQUIPMENT.equipment_model_id = EQUIPMENT_MODEL.equipment_model_id AND " +
			//                "       EQUIPMENT_MODEL.company_id = COMPANY.company_id AND " +
			//                "       EQUIPMENT_LOCATION.installation_date < CURRENT_TIMESTAMP AND " +
			//                "       EQUIPMENT_LOCATION.removal_date > CURRENT_TIMESTAMP AND " +
			//                "       SEISMIC_STREAM.start_time_stamp < CURRENT_TIMESTAMP AND " +
			//                "       SEISMIC_STREAM.stop_time_stamp > CURRENT_TIMESTAMP AND " +
			//                "       EQUIPMENT_MODEL.model_nmbr is not NULL AND " +
			//
			//                "       EQUIPMENT.equipment_parent_id IS NULL AND " +
			//                "       SEISMIC_STATION.station_id IN ('TRAB', 'TROB') AND " +
			//                " EXISTS(SELECT * FROM equipment_location l, data_logger d, equipment e WHERE l.place_id = PLACE.place_id AND l.equipment_id = d.equipment_id AND l.equipment_id != EQUIPMENT.equipment_id AND e.equipment_id = l.equipment_id AND e.equipment_model_id = EQUIPMENT.equipment_model_id AND l.installation_date < EQUIPMENT_LOCATION.removal_date AND l.removal_date > EQUIPMENT_LOCATION.installation_date) AND " +
			//                " PLACE.name NOT IN ('Avalon GeoNet Undeployed Equipment Store', 'GNS Wairakei Grange Building Storeroom', 'Wairakei GeoNet Undeployed Equipment Store') AND " +
			//                " EQUIPMENT_MODEL.model NOT IN ('EARSS/3', 'Taurus') AND " +
			//"       "+sensors+" "+
			"ORDER BY "+
			"       SEISMIC_STATION.station_id, SENSOR_LOCATION.installation_date"+
			"")
	*/

	for _, sql := range queries {
		log.Printf(sql)
		rows, err := db.Query(sql)

		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {

			var name string
			var lat float32
			var lon float32
			var sta string
			var loc string
			var cha string
			var net string
			var ext string
			var rate int32
			var gain float32

			if err := rows.Scan(&name, &lat, &lon, &sta, &loc, &cha, &net, &ext, &rate, &gain); err != nil {
				log.Fatal(err)
			}
			q, ok := Q[rate]
			if ok {
				var srcname string

				if ext != "" {
					srcname = fmt.Sprintf("%s_%s_%s_%s", ext, sta, loc, cha)
				} else {
					srcname = fmt.Sprintf("%s_%s_%s_%s", net, sta, loc, cha)
				}

				impacts[srcname] = Impact{name, (float32)(rate), gain, q, lat, lon}
			}
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
	}

	return impacts
}

func impact(w http.ResponseWriter, r *http.Request) {

	impacts := Impacts()
	fmt.Printf("%q\n", impacts)

	b, err := json.Marshal(impacts)
	if err != nil {
		log.Fatal(err)
	}

	io.WriteString(w, bytes.NewBuffer(b).String())
}
