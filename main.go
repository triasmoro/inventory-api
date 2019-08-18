package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/triasmoro/inventory-api/app"
)

func main() {
	port := 8080

	app, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	// generate tables
	if err := generate(app.DB); err != nil {
		log.Fatal(err)
	}
	log.Println("Database has been generated")

	// routing
	r := mux.NewRouter()

	srv := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, r),
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
