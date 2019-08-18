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

// PostPurchaseOrder to add new purchase order
// TODO: check uniqueness of po_number
func PostPurchaseOrder(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			WriteErrors(w, FieldErrors{{"read body", ErrFormatInvalid}})
			return
		}

		var order model.PurchaseOrder
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

		// check existing of product variant
		for _, detail := range order.Details {
			exist, err := store.IsProductVariantExist(detail.ProductVariantID)
			if err != nil {
				WriteErrors(w, FieldErrors{{"retrieve product variant", ErrFailed}})
				return
			} else if !exist {
				WriteErrors(w, FieldErrors{{"retrieve product variant", ErrNotFound}})
				return
			}
		}

		// save
		if err := store.SavePurchaseOrder(&order); err != nil {
			WriteErrors(w, FieldErrors{{"save order", ErrFailed}})
			return
		}

		WriteData(w, http.StatusOK, order)
	}
}
