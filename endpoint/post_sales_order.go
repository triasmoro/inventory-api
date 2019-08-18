package endpoint

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/triasmoro/inventory-api/app"
	"github.com/triasmoro/inventory-api/helper"
	"github.com/triasmoro/inventory-api/model"
)

func PostSalesOrder(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			WriteErrors(w, FieldErrors{{"read body", ErrFormatInvalid}})
			return
		}

		var order model.SalesOrder
		if err := json.Unmarshal(body, &order); err != nil {
			WriteErrors(w, FieldErrors{{"unmarshal body", ErrFormatInvalid}})
			return
		}

		store := app.Store

		// check time validity
		if _, err := time.Parse(helper.DateLayoutYMDHIS, order.Time); err != nil {
			WriteErrors(w, FieldErrors{{"time value", ErrFormatInvalid}})
			return
		}

		// generate so number
		orderNumber, err := helper.GenerateSONumber()
		if err != nil {
			WriteErrors(w, FieldErrors{{"generate order number", ErrFailed}})
			return
		}

		// assign new order number
		order.SONumber = orderNumber

		// save
		if err := store.SaveSalesOrder(&order); err != nil {
			WriteErrors(w, FieldErrors{{"save order", ErrFailed}})
			return
		}

		WriteData(w, http.StatusOK, order)
	}
}
