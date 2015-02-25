package main

import (
	"database/sql"
	_ "github.com/mattn/go-oci8"
        "log"
        "net/http"
        "os"
        "time"
)

var (
        db *sql.DB
        client *http.Client
)

func main() {
        var err error
        // TODO
        var port string = "9999"

        dsn := os.Getenv("DSN")
        if dsn == "" {
                log.Println("DSN environment variable not set, of the form \"\".")
                log.Fatal(err)
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

        timeout := time.Duration(5 * time.Second)
        client = &http.Client{
                Timeout: timeout,
        }

        http.HandleFunc("/places", places)

        log.Fatal(http.ListenAndServe(":"+port, nil))
}
