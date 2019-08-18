package endpoint

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/triasmoro/inventory-api/app"
	"github.com/triasmoro/inventory-api/helper"
	"github.com/triasmoro/inventory-api/model"
)

// notes availability
const (
	Lost     = "Barang Hilang"
	Damaged  = "Barang Rusak"
	Sampling = "Sample Barang"
)

// PostStockOut endpoint
func PostStockOut(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			WriteErrors(w, FieldErrors{{"read body", ErrFormatInvalid}})
			return
		}

		var stockOut model.StockOut
		if err := json.Unmarshal(body, &stockOut); err != nil {
			WriteErrors(w, FieldErrors{{"unmarshal body", ErrFormatInvalid}})
			return
		}

		store := app.Store

		// check time validity
		if stockOut.Time != "" {
			if _, err := time.Parse(helper.DateLayoutYMDHIS, stockOut.Time); err != nil {
				WriteErrors(w, FieldErrors{{"time value", ErrFormatInvalid}})
				return
			}
		}

		if stockOut.SalesOrderDetailID != 0 {
			// stock out based on sales order

			// check existing of sales order and get order number at once
			orderNumber, err := store.GetSONumberByOrderDetailID(stockOut.SalesOrderDetailID)
			if err != nil {
				WriteErrors(w, FieldErrors{{"retrieve sales order detail", ErrFailed}})
				return
			} else if orderNumber == "" {
				WriteErrors(w, FieldErrors{{"retrieve sales order detail", ErrNotFound}})
				return
			}

			// generate notes
			stockOut.Notes = fmt.Sprintf("Pesanan %s", orderNumber)

			// save
			if err := store.SaveStockOutBasedSalesOrder(&stockOut); err != nil {
				WriteErrors(w, FieldErrors{{"save stock out", ErrFailed}})
				return
			}

		} else if stockOut.ProductVariantID != 0 {
			// stock out based without sales order such as lost / damaged / sampling

			// check notes validity
			if stockOut.Notes != Lost &&
				stockOut.Notes != Damaged &&
				stockOut.Notes != Sampling {
				msg := fmt.Sprintf("notes (only `%s`, `%s`, `%s`)", Lost, Damaged, Sampling)
				WriteErrors(w, FieldErrors{{msg, ErrValueInvalid}})
				return
			}

			// check existing of product
			exist, err := store.IsProductVariantExist(stockOut.ProductVariantID)
			if err != nil {
				WriteErrors(w, FieldErrors{{"retrieve product variant", ErrFailed}})
				return
			} else if !exist {
				WriteErrors(w, FieldErrors{{"retrieve product variant", ErrNotFound}})
				return
			}

			// save
			if err := store.SaveStockOutBasedProduct(&stockOut); err != nil {
				WriteErrors(w, FieldErrors{{"save stock out", ErrFailed}})
				return
			}
		} else {
			WriteErrors(w, FieldErrors{{"sales order or product", ErrNotFound}})
			return
		}

		WriteData(w, http.StatusOK, stockOut)
	}
}
