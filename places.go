package main

import (
	"bytes"
	"encoding/json"
	_ "github.com/mattn/go-oci8"
	"io"
	"log"
	"net/http"
)

type Place struct {
	Name      string
	Latitude  float32
	Longitude float32
}

func places(w http.ResponseWriter, r *http.Request) {
	//var result string = ""
	rows, err := db.Query("SELECT name, latitude, longitude FROM place ORDER by name")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var places map[string]Place

	places = make(map[string]Place)
	for rows.Next() {
		var name string
		var lat, lon float32
		if err := rows.Scan(&name, &lat, &lon); err != nil {
			log.Fatal(err)
		}
		places[name] = Place{name, lat, lon}
		//result = result + fmt.Sprintf(" %s %g %g", name, lat, lon)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(places)
	if err != nil {
		log.Fatal(err)
	}

	io.WriteString(w, bytes.NewBuffer(b).String())
}
