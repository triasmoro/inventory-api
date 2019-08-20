package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/triasmoro/inventory-api/app"
	"github.com/triasmoro/inventory-api/route"
)

func main() {
	// assign port
	var port int
	flag.IntVar(&port, "p", 8080, "port")
	flag.Parse()

	app, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	// generate tables
	if err := generate(app.DB); err != nil {
		log.Fatal(err)
	}
	log.Println("Database has been generated")

	srv := &http.Server{
		Handler:      route.Router(app), // routes
		Addr:         fmt.Sprintf("127.0.0.1:%d", port),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Printf("Listening on port %d", port)
	log.Fatal(srv.ListenAndServe())
}

func generate(db *sql.DB) error {
	file, err := ioutil.ReadFile("db-sqlite.sql")
	if err != nil {
		return err
	}

	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		if _, err := db.Exec(request); err != nil {
			return err
		}
	}

	return nil
}
