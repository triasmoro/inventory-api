package main

import (
	"database/sql"
	"flag"
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
	"github.com/triasmoro/inventory-api/endpoint"
)

func main() {
	// config port
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

	// routing
	r := mux.NewRouter()
	r.HandleFunc("/product", endpoint.PostProduct(app)).Methods("POST").Headers("Content-Type", "application/json")
	r.HandleFunc("/product/{id:[0-9]+}", endpoint.PatchProduct(app)).Methods("PATCH")
	r.HandleFunc("/product_variant/{id:[0-9]+}", endpoint.DeleteProductVariant(app)).Methods("DELETE")
	r.HandleFunc("/purchase_order", endpoint.PostPurchaseOrder(app)).Methods("POST").Headers("Content-Type", "application/json")
	r.HandleFunc("/purchase_order/{id:[0-9]+}", endpoint.DeletePurchaseOrder(app)).Methods("DELETE")
	r.HandleFunc("/stock_in", endpoint.PostStockIn(app)).Methods("POST").Headers("Content-Type", "application/json")
	r.HandleFunc("/stock_in/{id:[0-9]+}", endpoint.DeleteStockIn(app)).Methods("DELETE")
	r.HandleFunc("/sales_order", endpoint.PostSalesOrder(app)).Methods("POST").Headers("Content-Type", "application/json")
	r.HandleFunc("/sales_order/{id:[0-9]+}", endpoint.DeleteSalesOrder(app)).Methods("DELETE")
	r.HandleFunc("/stock_out", endpoint.PostStockOut(app)).Methods("POST").Headers("Content-Type", "application/json")
	r.HandleFunc("/stock_out/{id:[0-9]+}", endpoint.DeleteStockOut(app)).Methods("DELETE")
	r.HandleFunc("/export/product", endpoint.GetExportProduct(app)).Methods("GET")

	// report
	r.HandleFunc("/actual_stock", endpoint.GetActualStock(app)).Methods("GET")
	r.HandleFunc("/assets_report", endpoint.GetAssetsReport(app)).Methods("GET")
	r.HandleFunc("/sales_report", endpoint.GetSalesReport(app)).Methods("GET")

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
