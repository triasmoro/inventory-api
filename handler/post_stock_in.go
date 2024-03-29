package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/triasmoro/inventory-api/app"
	"github.com/triasmoro/inventory-api/helper"
	"github.com/triasmoro/inventory-api/model"
)

// PostStockIn endpoint
func PostStockIn(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			WriteErrors(w, FieldErrors{{"read body", ErrFormatInvalid}})
			return
		}

		var stockIn model.StockIn
		if err := json.Unmarshal(body, &stockIn); err != nil {
			WriteErrors(w, FieldErrors{{"unmarshal body", ErrFormatInvalid}})
			return
		}

		store := app.Store

		// check time validity
		if _, err := time.Parse(helper.DateLayoutYMDHIS, stockIn.Time); err != nil {
			WriteErrors(w, FieldErrors{{"time value", ErrFormatInvalid}})
			return
		}

		// check existing of purchase order detail
		exist, err := store.IsPurchaseOrderDetailExist(stockIn.PurchaseOrderDetailID)
		if err != nil {
			WriteErrors(w, FieldErrors{{"retrieve purchase order detail", ErrFailed}})
			return
		} else if !exist {
			WriteErrors(w, FieldErrors{{"retrieve purchase order detail", ErrNotFound}})
			return
		}

		// save
		if err := store.SaveStockIn(&stockIn); err != nil {
			WriteErrors(w, FieldErrors{{"save stock in", ErrFailed}})
			return
		}

		WriteData(w, http.StatusOK, stockIn)
	}
}
