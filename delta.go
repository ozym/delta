package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	db     *sql.DB
	client *http.Client
)

func main() {
	var err error

	// runtime settings
	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "make noise")
	var server bool
	flag.BoolVar(&server, "server", false, "run as a web service")
	var port string = "9999"
	flag.StringVar(&port, "port", "9999", "port to run the service on")
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 5.0*time.Second, "service timeout")

	var do_impact bool
	flag.BoolVar(&do_impact, "impact", false, "output impact data")

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
		impacts := Impacts()

		b, err := json.Marshal(impacts)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(bytes.NewBuffer(b).String())
	}

	if server {
		client = &http.Client{
			Timeout: timeout,
		}

		http.HandleFunc("/impact", impact)

		log.Fatal(http.ListenAndServe(":"+port, nil))
	}

}
