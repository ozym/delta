package main

import (
	"encoding/json"
	_ "github.com/mattn/go-oci8"
)

type Place struct {
	Name      string  `json:"Name"`
	Latitude  float32 `json:"Latitude"`
	Longitude float32 `json:"Longitude"`
}

func PlacesConfig(indent bool) ([]byte, error) {
	places, err := Places()
	if err != nil {
		return nil, err
	}

	if indent {
		return json.MarshalIndent(places, "", "  ")
	} else {
		return json.Marshal(places)
	}
}

func Places() (map[string]Place, error) {
	var places map[string]Place

	places = make(map[string]Place)

	rows, err := db.Query("SELECT name, latitude, longitude FROM place ORDER by name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var lat, lon float32
		if err := rows.Scan(&name, &lat, &lon); err != nil {
			return nil, err
		}
		places[name] = Place{name, lat, lon}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return places, nil
}
