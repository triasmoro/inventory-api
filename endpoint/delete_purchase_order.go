package endpoint

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/triasmoro/inventory-api/app"
)

// DeletePurchaseOrder endpoint
func DeletePurchaseOrder(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])

		store := app.Store

		// check stock-in from this po
		exist, err := store.IsStockInExistByPurchaseOrderID(id)
		if err != nil {
			WriteErrors(w, FieldErrors{{"retrieve stock-in", ErrFailed}})
			return
		} else if exist {
			WriteErrors(w, FieldErrors{{"stock-in exist", ErrFailed}})
			return
		}

		if err := store.DeletePurchaseOrder(id); err != nil {
			WriteErrors(w, FieldErrors{{"delete purchase order", ErrFailed}})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
