package main

import (
	"encoding/json"
	_ "github.com/mattn/go-oci8"
)

type Station struct {
	Code      string  `json:"Code"`
	Name      string  `json:"Name"`
	Network   string  `json:"Network"`
	External  *string `json:"External,omitempty"`
	Latitude  float32 `json:"Latitude"`
	Longitude float32 `json:"Longitude"`
	Height    float32 `json:"Height"`
}

func StationsConfig(indent bool) ([]byte, error) {
	stations := make(map[string]Station)
	list, err := Stations()
	if err != nil {
		return nil, err
	}
	for _, s := range list {
		stations[s.Code] = s
	}

	if indent {
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
		if err := rows.Scan(&s.Code, &s.Name, &s.Latitude, &s.Longitude, &s.Height, &s.Network, &s.External); err != nil {
			return nil, err
		}
		if s.External == nil {
			s.External = &s.Network
		}
		stations = append(stations, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stations, nil
}
