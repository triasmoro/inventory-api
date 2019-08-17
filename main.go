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
)

func main() {
	port := 8080

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// generate tables
	generate(db)

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

func generate(db *sql.DB) {
	file, err := ioutil.ReadFile("db.sql")
	if err != nil {
		log.Fatal(err)
	}

	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		if _, err := db.Exec(request); err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Database has been generated")
}
