package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"log"
	"os"
)

var db *sql.DB

var CSD = []string{"CUSP3D SENSOR"}

func store(output string, config []byte) error {
	if output != "-" && output != "" {
		f, err := os.Create(output)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(config)
		if err != nil {
			return err
		}
	} else {
		fmt.Print(bytes.NewBuffer(config).String())
	}
	return nil
}

func main() {
	var err error

	// runtime settings
	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "make noise")

	// where to store results
	var output string
	flag.StringVar(&output, "output", "-", "store to run the service on")
	var indent bool
	flag.BoolVar(&indent, "pretty", false, "produce indented json output")

	// output config files
	var do_impact bool
	flag.BoolVar(&do_impact, "impact", false, "output impact data")
	var do_places bool
	flag.BoolVar(&do_places, "places", false, "output places data")
	var do_csd bool
	flag.BoolVar(&do_csd, "csd", false, "output csd data")

	// oracle connection details
	var dsn string
	flag.StringVar(&dsn, "dsn", "", "provide DSN connection string, overides env variable \"DSN\"")

	flag.Parse()

	if dsn == "" {
		dsn = os.Getenv("DSN")
		if dsn == "" {
			log.Fatal("DSN environment variable not set, of the form \"<user>/<password>@<server>:<port>/<instance>\".")
		}
	}

	db, err = sql.Open("oci8", dsn)
	if err != nil {
		log.Println("Unable to build SQL database connection")
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Println("Unable to connect to the database connection")
	}

	if do_impact {
		config, err := ImpactConfig(indent)
		if err != nil {
			log.Fatal(err)
		}
		err = store(output, config)
		if err != nil {
			log.Fatal(err)
		}
	}
	if do_places {
		config, err := PlacesConfig(indent)
		if err != nil {
			log.Fatal(err)
		}
		err = store(output, config)
		if err != nil {
			log.Fatal(err)
		}
	}

	/*
		m, err := GetEquipmentModels()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(m)

		mm, err := GetEquipmentModel(2)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(mm)

		s, err := SeismicStreams()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(s)
	*/

	if do_csd {
		config, err := CsdConfig(indent, CSD)
		if err != nil {
			log.Fatal(err)
		}
		err = store(output, config)
		if err != nil {
			log.Fatal(err)
		}
	}

}
