package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"reflect"
	"strings"
)

type Station struct {
	Code      string  `json:"Code" etcd:"code"`
	Name      string  `json:"Name" etcd:"name"`
	Network   string  `json:"Network" etcd:"network"`
	External  string  `json:"External" etcd:"external"`
	Latitude  float32 `json:"Latitude" etcd:"latitude"`
	Longitude float32 `json:"Longitude" etcd:"longitude"`
	Height    float32 `json:"Height" etcd:"height"`
}

func (s *Station) etcd() map[string]string {
	results := make(map[string]string)

	val := reflect.ValueOf(s).Elem()
	for i := 0; i < val.NumField(); i++ {
		t := val.Type().Field(i).Tag.Get("etcd")
		if t == "" {
			continue
		}
		results[t] = fmt.Sprintf("%v", val.Field(i).Interface())
		if strings.ContainsAny(results[t], " \t") {
			results[t] = "\"" + results[t] + "\""
		}
	}

	return results
}

func StationsConfig() ([]byte, error) {
	stations := make(map[string]Station)
	list, err := Stations()
	if err != nil {
		return nil, err
	}
	for _, s := range list {
		stations[s.Code] = s
	}

	if etcd {
		var results [][]byte
		for k, s := range stations {

			elements := s.etcd()
			for e, v := range elements {
				r := ([]byte)("/delta/seismic_station/" + strings.ToLower(k) + "/" + e + " " + v)
				results = append(results, r)
			}
		}
		results = append(results, []byte(""))
		return bytes.Join(results, ([]byte)("\n")), nil
	} else if indent {
		return json.MarshalIndent(stations, "", "  ")
	} else {
		return json.Marshal(stations)
	}
}

func Stations() ([]Station, error) {
	var stations []Station

	rows, err := db.Query("SELECT s.station_id, s.long_name, s.latitude, s.longitude, s.height, n.code, n.external_code FROM NETWORK n JOIN SEISMIC_STATION s ON s.network_id = n.network_id WHERE s.long_name is NOT NULL")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := Station{}
		var e *string
		if err := rows.Scan(&s.Code, &s.Name, &s.Latitude, &s.Longitude, &s.Height, &s.Network, &e); err != nil {
			return nil, err
		}
		s.External = s.Network
		if e != nil {
			s.External = *e
		}
		stations = append(stations, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stations, nil
}
